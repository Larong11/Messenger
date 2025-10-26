package user

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

func GeneratePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
func Generate6DigitCode() string {
	rand.Seed(time.Now().UnixNano()) // инициализация генератора случайных чисел
	code := rand.Intn(1000000)       // число от 0 до 999999
	return fmt.Sprintf("%06d", code) // всегда 6 цифр с ведущими нулями
}
