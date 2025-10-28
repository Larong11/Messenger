package services

import "context"

type EmailService interface {
	SendVerificationCode(ctx context.Context, to string, code string) error
}
