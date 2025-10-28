package user

import (
	"context"
	"errors"
	"server/domain/repositories"
	"server/domain/services"
	packageuser "server/domain/user"
	"time"
)

const CodeTTL = 10 * time.Minute

type RegisterUserUseCases struct {
	userRepo     repositories.UserRepository
	emailService services.EmailService
}

func NewRegisterUserUseCases(ur repositories.UserRepository, es services.EmailService) *RegisterUserUseCases {
	return &RegisterUserUseCases{
		userRepo:     ur,
		emailService: es,
	}
}
func (uc *RegisterUserUseCases) CheckUserName(ctx context.Context, userName string) (bool, error) {
	err := packageuser.ValidateUserName(userName)
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
	err := packageuser.ValidateEmail(email)
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

func (uc *RegisterUserUseCases) RegisterUser(ctx context.Context, firstName, lastName, userName, email, password, avatarURL, verificationCode string) (*int, error) {
	err := packageuser.ValidateUserName(userName)
	if err != nil {
		return nil, err
	}
	checkID, err := uc.userRepo.FindByUserName(ctx, userName)
	if err != nil {
		return nil, err
	}
	if checkID != nil {
		return nil, errors.New("user with such username already exists")
	}

	err = packageuser.ValidateEmail(email)
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

	passwordHash, err := packageuser.GeneratePasswordHash(password)
	if err != nil {
		return nil, err
	}
	user, err := packageuser.NewUser(firstName, lastName, userName, email, passwordHash, avatarURL)
	if err != nil {
		return nil, err
	}

	storedVerificationCode, createdAt, err := uc.userRepo.GetVerificationCode(ctx, email)
	if err != nil {
		return nil, err
	}
	if time.Since(createdAt) > CodeTTL {
		return nil, errors.New("verification code expired")
	}
	if storedVerificationCode != verificationCode {
		return nil, errors.New("incorrect verification code")
	}
	ID, err := uc.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return &ID, nil
}
func (uc *RegisterUserUseCases) RequestVerificationCode(ctx context.Context, email string) error {
	err := packageuser.ValidateEmail(email)
	if err != nil {
		return err
	}
	ID, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	if ID != nil {
		return errors.New("user with such email exists")
	}
	code, err := packageuser.Generate6DigitCode()
	if err != nil {
		return err
	}
	err = uc.userRepo.CreateVerificationCode(ctx, email, code)
	if err != nil {
		return err
	}
	//err = uc.emailService.SendVerificationCode(ctx, email, code)// TODO функция пока не рабочая, нет открытых портов
	//if err != nil {
	//	return err
	//}
	return nil
}
