package dto

type OtpDto struct {
	RetryAttemptsLeft int64
	ExpiresIn         int64
	RetryAfterIn      int64
	IsValid           bool
}
