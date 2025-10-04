package entities

type OTP struct {
	Number string `json:"otp_number"`
}

type OTPService struct {
	OTPNumber         string
	Expiration        int64
	MaxAttempt        int64
	RetryInterval     int64
	RetryAttemptsLeft int64
}
