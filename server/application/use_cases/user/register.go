package user

import (
	"context"
	"errors"
	"server/domain/repositories"
	package_user "server/domain/user"
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
	err := package_user.ValidateUserName(userName)
	if err != nil {
		return false, err
	}
	ID, err := uc.userRepo.FindByUserName(ctx, userName)
	if err != nil {
		return false, err
	}
	if ID != nil {
		return false, nil
	}
	return true, nil
}
func (uc *RegisterUserUseCases) CheckEmail(ctx context.Context, email string) (bool, error) {
	err := package_user.ValidateEmail(email)
	if err != nil {
		return false, err
	}
	ID, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	if ID != nil {
		return false, nil
	}
	return true, nil
}
func (uc *RegisterUserUseCases) RegisterUser(ctx context.Context, firstName, lastName, userName, email, password, avatarURL string) (*int, error) {
	err := package_user.ValidateUserName(userName)
	if err != nil {
		return nil, err
	}
	checkID, err := uc.userRepo.FindByUserName(ctx, userName)
	if err != nil {
		return nil, err
	}
	if checkID != nil {
		return nil, errors.New("username with such username already exists")
	}

	err = package_user.ValidateEmail(email)
	if err != nil {
		return nil, err
	}

	checkID, err = uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if checkID != nil {
		return nil, errors.New("user with such email already exists")
	}

	passwordHash, err := package_user.GeneratePasswordHash(password)
	if err != nil {
		return nil, err
	}
	user, err := package_user.NewUser(firstName, lastName, userName, email, passwordHash, avatarURL)
	if err != nil {
		return nil, err
	}
	ID, err := uc.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return &ID, nil
}
