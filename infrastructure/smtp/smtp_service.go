package smtp

import (
	"fmt"

	"github.com/winartodev/apollo-be/config"
	"gopkg.in/gomail.v2"
)

// SMTPService defines the interface for email sending functionality
type SMTPService interface {
	SendText(recipient string, subject string, body string) error
	SendHTML(recipient string, subject string, htmlBody string) error
}

// smtpService implements SMTPService
type smtpService struct {
	host     string
	port     int
	sender   string
	password string
}

// NewSMTPService creates a new SMTP service instance
func NewSMTPService(smtpConfig *config.SMTPConfig) SMTPService {
	return &smtpService{
		host:     smtpConfig.Host,
		port:     smtpConfig.Port,
		sender:   smtpConfig.Sender,
		password: smtpConfig.Password,
	}
}

// SendText sends a plain text email
func (s *smtpService) SendText(recipient string, subject string, body string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", s.sender)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(s.host, s.port, s.sender, s.password)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send text email: %w", err)
	}

	return nil
}

// SendHTML sends an HTML email
func (s *smtpService) SendHTML(recipient string, subject string, htmlBody string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.sender)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(s.host, s.port, s.sender, s.password)
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send HTML email: %w", err)
	}

	return nil
}
