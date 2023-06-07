package otp

import (
	"encoding/base32"
	"fmt"
	"net/url"
	"strings"
	"time"
)

const (
	OtpTypeTotp = "totp"
	OtpTypeHotp = "hotp"
)

// TOTP time-based OTP counters.
type TOTP struct {
	OTP
	interval int
}

func NewTOTP(secret string, digits, interval int, hasher *Hasher) *TOTP {
	otp := NewOTP(secret, digits, hasher)
	return &TOTP{OTP: otp, interval: interval}
}

func GenerateTOTP(config Config) OTPResponse {
	if config.Digits == 0 {
		config.Digits = DefaultDigits
	}
	if config.Interval == 0 {
		config.Interval = DefaultInterval
	}
	bt := base32.StdEncoding.EncodeToString([]byte(config.Secret))
	totp := NewTOTP(bt, config.Digits, config.Interval, nil)
	return totp.NowWithExpiration()
}

func VerifyTOTP(secret string, otp string, timestamp int) bool {
	bt := base32.StdEncoding.EncodeToString([]byte(secret))
	totp := NewTOTP(bt, DefaultDigits, DefaultInterval, nil)
	return totp.Verify(otp, timestamp)
}

// At Generate time OTP of given timestamp
func (t *TOTP) At(timestamp int) string {
	return t.generateOTP(t.timecode(timestamp))
}

// Now Generate the current time OTP
func (t *TOTP) Now() string {
	return t.At(CurrentTimestamp())
}

// NowWithExpiration Generate the current time OTP and expiration time
func (t *TOTP) NowWithExpiration() OTPResponse {
	interval64 := int64(t.interval)
	timeCodeInt64 := time.Now().Unix() / interval64
	expirationTime := (timeCodeInt64 + 1) * interval64
	expireTime := time.Unix(expirationTime, 0)
	return OTPResponse{
		Code:           t.generateOTP(int(timeCodeInt64)),
		ExpirationTime: expirationTime,
		ExpiresIn:      int64(expireTime.Sub(time.Now()).Seconds()),
	}
}

/*
Verify OTP.
params:

	otp:         the OTP to check against
	timestamp:   time to check OTP at
*/
func (t *TOTP) Verify(otp string, timestamp int) bool {
	return otp == t.At(timestamp)
}

// ProvisioningURI /*
func (t *TOTP) ProvisioningURI(accountName, issuerName string) string {
	return BuildUri(
		OtpTypeTotp,
		t.secret,
		accountName,
		issuerName,
		t.hasher.HashName,
		0,
		t.digits,
		t.interval)
}

func (t *TOTP) timecode(timestamp int) int {
	return timestamp / t.interval
}

// BuildUri /*
func BuildUri(otpType, secret, accountName, issuerName, algorithm string, initialCount, digits, period int) string {
	if otpType != OtpTypeHotp && otpType != OtpTypeTotp {
		panic("otp type error, got " + otpType)
	}

	urlParams := make([]string, 0)
	urlParams = append(urlParams, "secret="+secret)
	if otpType == OtpTypeHotp {
		urlParams = append(urlParams, fmt.Sprintf("counter=%d", initialCount))
	}
	label := url.QueryEscape(accountName)
	if issuerName != "" {
		issuerNameEscape := url.QueryEscape(issuerName)
		label = issuerNameEscape + ":" + label
		urlParams = append(urlParams, "issuer="+issuerNameEscape)
	}
	if algorithm != "" && algorithm != "sha1" {
		urlParams = append(urlParams, "algorithm="+strings.ToUpper(algorithm))
	}
	if digits != 0 && digits != 6 {
		urlParams = append(urlParams, fmt.Sprintf("digits=%d", digits))
	}
	if period != 0 && period != 30 {
		urlParams = append(urlParams, fmt.Sprintf("period=%d", period))
	}
	return fmt.Sprintf("otpauth://%s/%s?%s", otpType, label, strings.Join(urlParams, "&"))
}

// CurrentTimestamp get current timestamp
func CurrentTimestamp() int {
	return int(time.Now().Unix())
}

// Itob integer to byte array
func Itob(integer int) []byte {
	byteArr := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		byteArr[i] = byte(integer & 0xff)
		integer = integer >> 8
	}
	return byteArr
}
