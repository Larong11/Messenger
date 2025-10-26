package user

import (
	"errors"
	"net/mail"
	"regexp"
)

func ValidateUserName(userName string) error {
	if len(userName) < 5 || len(userName) > 20 {
		return errors.New("username must be 3–20 characters long")
	}
	valid := regexp.MustCompile(`^[A-Za-z0-9_]+$`)
	if !valid.MatchString(userName) {
		return errors.New("username may contain only letters, digits, and underscores")
	}
	return nil
}
func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("invalid email format")
	}
	return nil
}
