package user

import (
	"net/mail"
	"regexp"
	"server/internal/errors"
)

func ValidateUserName(userName string) error {
	if len(userName) < 5 || len(userName) > 20 {
		return errors.NewBadRequest("username must be 3–20 characters long")
	}
	valid := regexp.MustCompile(`^[A-Za-z0-9_]+$`)
	if !valid.MatchString(userName) {
		return errors.NewBadRequest("username may contain only letters, digits, and underscores")
	}
	return nil
}
func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.NewBadRequest("invalid email format")
	}
	if len(email) > 255 {
		return errors.NewBadRequest("email too long")
	}
	return nil
}
