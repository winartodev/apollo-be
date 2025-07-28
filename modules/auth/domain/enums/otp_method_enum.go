package enums

type OtpMethod int64

const (
	Email OtpMethod = iota
	SMS
)

func (s OtpMethod) String() string {
	switch s {
	case Email:
		return "email"
	case SMS:
		return "sms"
	}
	return "unknown"
}
