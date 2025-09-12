package helper

import "gopkg.in/gomail.v2"

type SmtpConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Sender   string `yaml:"sender"`
	Password string `yaml:"password"`
}

func (h *SmtpConfig) SendSmtp(recipient string, subject string, body string) (err error) {
	m := gomail.NewMessage()

	m.SetHeader("From", h.Sender)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(h.Host, h.Port, h.Sender, h.Password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func (h *SmtpConfig) SendSmtpHTML(recipient string, subject string, htmlBody string) (err error) {
	m := gomail.NewMessage()
	m.SetHeader("From", h.Sender)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(h.Host, h.Port, h.Sender, h.Password)
	return d.DialAndSend(m)
}
