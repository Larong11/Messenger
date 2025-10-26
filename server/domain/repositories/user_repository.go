package repositories

import (
	"context"
	"server/domain/user"
)

type UserRepository interface {
	FindByUserName(ctx context.Context, userName string) (*int, error)
	FindByEmail(ctx context.Context, email string) (*int, error)
	CreateUserWithVerificationCode(ctx context.Context, user *user.User, verificationCode string) (int, error)
}
