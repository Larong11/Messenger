package repositories

import "context"

type UserRepository interface {
	FindByUserName(ctx context.Context, userName string) (*int, error)
	FindByEmail(ctx context.Context, email string) (*int, error)
}
