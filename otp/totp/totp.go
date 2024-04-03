/**
 *  Copyright 2014 Paul Querna
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package totp

import (
	"encoding/base32"
	"io"
	"math"
	"net/url"
	"strconv"
	"time"

	"github.com/oarkflow/pkg/otp"
	"github.com/oarkflow/pkg/otp/hotp"
)

// ValidateOpts provides options for ValidateCustom().
type ValidateOpts struct {
	// Number of seconds a TOTP hash is valid for. Defaults to 30 seconds.
	Period uint
	// Periods before or after the current time to allow.  Value of 1 allows up to Period
	// of either side of the specified time.  Defaults to 0 allowed skews.  Values greater
	// than 1 are likely sketchy.
	Skew uint
	// Digits as part of the input. Defaults to 6.
	Digits otp.Digits
	// Algorithm to use for HMAC. Defaults to SHA1.
	Algorithm otp.Algorithm
	// the time in which we would like to validate our code
	// in the normal usage, it is equal to current time : time.Now()
	// but for testing puposes, it could be changed to a later/future time
	t time.Time
}

// GenerateOpts provides options for Generate().  The default values
// are compatible with Google-Authenticator.
type GenerateOpts struct {
	// Name of the issuing Organization/Company.
	Issuer string
	// Name of the User's Account (eg, email address)
	AccountName string
	// Number of seconds a TOTP hash is valid for. Defaults to 30 seconds.
	Period uint
	// Size in size of the generated Secret. Defaults to 20 bytes.
	SecretSize uint
	// Secret to store. Defaults to a randomly generated secret of SecretSize.  You should generally leave this empty.
	Secret []byte
	// Digits to request. Defaults to 6.
	Digits otp.Digits
	// Algorithm to use for HMAC. Defaults to SHA1.
	Algorithm otp.Algorithm
	// Reader to use for generating TOTP Key.
	Rand io.Reader
}

var b32NoPadding = base32.StdEncoding.WithPadding(base32.NoPadding)

// ValidateWithOpts validate with opts
// it has deprecated Validate and will replace it soon.
func ValidateWithOpts(passcode, secret string, validateOpts ...ValidateOpt) (bool, error) {
	return validateCustomOpt(passcode, secret, validateOpts...)
}

// validateCustomOpt validates a TOTP given a user specified time and custom options.
// Most users should use Validate() to provide an interpolatable TOTP experience.
// This replicates ValidateCustomOpt
func validateCustomOpt(passcode, secret string, validateOpts ...ValidateOpt) (bool, error) {

	opts := new(ValidateOpts)

	for _, opt := range validateOpts {
		opt(opts)
	}
	opts.defaultOpts()

	counters := []uint64{}
	counter := int64(math.Floor(float64(opts.t.Unix()) / float64(opts.Period)))

	counters = append(counters, uint64(counter))
	for i := 1; i <= int(opts.Skew); i++ {
		counters = append(counters, uint64(counter+int64(i)))
		counters = append(counters, uint64(counter-int64(i)))
	}

	for _, counter := range counters {
		rv, err := hotp.ValidateCustom(passcode, counter, secret, hotp.ValidateOpts{
			Digits:    opts.Digits,
			Algorithm: opts.Algorithm,
		})

		if err != nil {
			return false, err
		}

		if rv {
			return true, nil
		}
	}

	return false, nil
}

// GenerateWithOpts a new TOTP Key.
// This func generates *otp.Key with desired GenerateOpts
// Example usage:
// GenerateWitOpts(WithAccountName("example account"))
func GenerateWithOpts(genOpts ...GenerateOpt) (*otp.Key, error) {

	opts := new(GenerateOpts)

	for _, opt := range genOpts {
		opt(opts)
	}

	if err := opts.defaults(); err != nil {
		return nil, err
	}
	// otpauth://totp/Example:alice@google.com?secret=JBSWY3DPEHPK3PXP&issuer=Example

	v := url.Values{}

	if opts.Secret != nil {
		v.Set("secret", b32NoPadding.EncodeToString(opts.Secret))
	} else {
		secret := make([]byte, opts.SecretSize)
		_, err := opts.Rand.Read(secret)
		if err != nil {
			return nil, err
		}
		v.Set("secret", b32NoPadding.EncodeToString(secret))
	}

	v.Set("issuer", opts.Issuer)
	v.Set("period", strconv.FormatUint(uint64(opts.Period), 10))
	v.Set("algorithm", opts.Algorithm.String())
	v.Set("digits", opts.Digits.String())

	u := url.URL{
		Scheme:   "otpauth",
		Host:     "totp",
		Path:     "/" + opts.Issuer + ":" + opts.AccountName,
		RawQuery: v.Encode(),
	}
	return otp.NewKeyFromURL(u.String())
}

// GenerateCodeWithOpts takes a time point and produces a passcode using a
// secret and the provided opts. (Under the hood, this is making an adapted
// call to hotp.GenerateCodeCustom)
func GenerateCodeWithOpts(secret string, validateOpts ...ValidateOpt) (passcode string, err error) {

	opts := new(ValidateOpts)

	for _, opt := range validateOpts {
		opt(opts)
	}
	opts.defaultOpts()

	counter := uint64(math.Floor(float64(opts.t.Unix()) / float64(opts.Period)))
	passcode, err = hotp.GenerateCodeCustom(secret, counter, hotp.ValidateOpts{
		Digits:    opts.Digits,
		Algorithm: opts.Algorithm,
	})
	if err != nil {
		return "", err
	}
	return passcode, nil
}
