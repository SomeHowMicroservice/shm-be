package smtp

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/smtp"
)

type mailerImpl struct {
	cfg  *MailerConfig
	auth smtp.Auth
}

//go:embed template/auth.html
var emailTemplates embed.FS

func NewMailer(cfg *MailerConfig) Mailer {
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	return &mailerImpl{
		cfg:  cfg,
		auth: auth,
	}
}

func (s *mailerImpl) Send(to, subject, body string) error {
	msg := []byte(fmt.Sprintf("Subject: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s", subject, body))
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	return smtp.SendMail(addr, s.auth, s.cfg.Username, []string{to}, msg)
}

func (s *mailerImpl) SendAuthEmail(to, subject, otp string) error {
	tmpl, err := template.ParseFS(emailTemplates, "template/auth.html")
	if err != nil {
		return fmt.Errorf("không thể load template: %w", err)
	}

	var body bytes.Buffer
	data := struct {
		Subject string 
		Otp string
	} {
		Subject: subject,
		Otp: otp,
	}

	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("không thể render template: %w", err)
	}

	return s.Send(to, subject, body.String())
}