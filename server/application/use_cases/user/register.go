package user

import (
	"context"
	"server/domain/repositories"
	"server/domain/user"
)

type RegisterUserUseCases struct {
	userRepo repositories.UserRepository
}

func NewRegisterUserUseCases(ur repositories.UserRepository) *RegisterUserUseCases {
	return &RegisterUserUseCases{
		userRepo: ur,
	}
}
func (uc *RegisterUserUseCases) CheckUserName(ctx context.Context, userName string) (bool, error) {
	err := user.ValidateUserName(userName)
	if err != nil {
		return false, err
	}
	ID, err := uc.userRepo.FindByUserName(ctx, userName)
	if err != nil {
		return false, err
	}
	if ID == nil {
		return false, nil
	}
	return true, nil
}
func (uc *RegisterUserUseCases) CheckEmail(ctx context.Context, email string) (bool, error) {
	err := user.ValidateEmail(email)
	if err != nil {
		return false, err
	}
	ID, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	if ID == nil {
		return false, nil
	}
	return true, nil
}
