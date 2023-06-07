package otp

const (
	DefaultAlgorithm = "sha1"
	DefaultDigits    = 6
	// DefaultInterval is the period parameter defines a period that a TOTP code will be valid for, in seconds.
	// ref: https://github.com/google/google-authenticator/wiki/Key-Uri-Format#period
	DefaultInterval = 300
)

type Parameter struct {
	Secret  string
	Issuer  string
	Account string
	UserID  uint
}

type Config struct {
	Secret   string `json:"secret"`
	Digits   int    `json:"digits"`
	Interval int    `json:"interval"`
}
