package enums

import "fmt"

type OtpOperationEnum string

const (
	OtpSignUp       OtpOperationEnum = "signup"
	OtpRequestReset OtpOperationEnum = "request_reset"
)

func ParseOtpOperationEnum(s string) (OtpOperationEnum, error) {
	switch s {
	case string(OtpSignUp):
		return OtpSignUp, nil
	case string(OtpRequestReset):
		return OtpRequestReset, nil
	default:
		return "", fmt.Errorf("invalid OtpOperationEnum: %s", s)
	}
}
