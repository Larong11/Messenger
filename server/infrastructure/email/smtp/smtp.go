package smtp

import (
	"context"
	"fmt"
	"net/smtp"
	"server/internal/errors"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

type EmailService struct {
	config Config
}

func NewEmailService(config Config) *EmailService {
	return &EmailService{config: config}
}

func (s *EmailService) SendVerificationCode(ctx context.Context, to string, code string) error {
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	subject := "Email Verification Code"
	body := fmt.Sprintf("Your verification code is: %s\n\nThis code will expire in 10 minutes.", code)

	message := []byte(
		"From: " + s.config.From + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"utf-8\"\r\n" +
			"\r\n" +
			body + "\r\n")

	addr := s.config.Host + ":" + s.config.Port
	err := smtp.SendMail(addr, auth, s.config.From, []string{to}, message)
	if err != nil {
		return errors.NewInternal("failed to send email")
	}
	return nil
}
