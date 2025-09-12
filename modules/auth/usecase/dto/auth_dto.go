package dto

type AuthDto struct {
	AccessToken  string
	RefreshToken string
	Otp          *OtpDto
}
