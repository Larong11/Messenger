package repositories

import (
	"context"
	"server/domain/user"
	"time"
)

type UserRepository interface {
	FindByUserName(ctx context.Context, userName string) (*int, error)
	FindByEmail(ctx context.Context, email string) (*int, error)
	CreateUser(ctx context.Context, user *user.User) (int, error)
	CreateVerificationCode(ctx context.Context, email string, verificationCode string) error
	GetVerificationCode(ctx context.Context, email string) (string, time.Time, error)
}
