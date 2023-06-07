package otp

// HOTP HMAC-based OTP counters.
type HOTP struct {
	OTP
}

func NewHOTP(secret string, digits int, hasher *Hasher) *HOTP {
	otp := NewOTP(secret, digits, hasher)
	return &HOTP{OTP: otp}

}

func GenerateHOTP(secret string) *HOTP {
	return NewHOTP(secret, 6, nil)
}

// At Generates the OTP for the given count.
func (h *HOTP) At(count int) string {
	return h.generateOTP(count)
}

/*
Verify OTP.
params:

	otp:   the OTP to check against
	count: the OTP HMAC counter
*/
func (h *HOTP) Verify(otp string, count int) bool {
	return otp == h.At(count)
}

// ProvisioningURI /*
func (h *HOTP) ProvisioningURI(accountName, issuerName string, initialCount int) string {
	return BuildUri(
		OtpTypeHotp,
		h.secret,
		accountName,
		issuerName,
		h.hasher.HashName,
		initialCount,
		h.digits,
		0)
}
