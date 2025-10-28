package user

import (
	"crypto/rand"
	"server/internal/errors"

	"golang.org/x/crypto/bcrypt"

	"math/big"
)

func GeneratePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.NewInternal("failed to generate password hash")
	}
	return string(hash), nil
}

func Generate6DigitCode() (string, error) {
	const digits = "0123456789"
	code := make([]byte, 6)

	for i := 0; i < 6; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", errors.NewInternal("failed to generate 6 digit code")
		}
		code[i] = digits[num.Int64()]
	}

	return string(code), nil
}
