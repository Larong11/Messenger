package user

import (
	"server/internal/errors"
	"time"
)

type Status int // симуляция enum

const (
	Offline Status = iota // Отчет начинается с 0
	Online
)

type User struct {
	ID           int    // string or int зависит от того будет ли генерировать его сервер или бд
	FirstName    string // длина не больше 50
	LastName     string // длина не больше 50
	UserName     string // Буквы английского алфавита и _, должен быть уникальным в бд длина не больше 50
	Email        string // Длина не больше 100
	PasswordHash string //длина не больше 255
	CreatedAt    time.Time
	AvatarURL    string
	LastSeenAt   time.Time
	UserStatus   Status
}

func NewUser(firstName, lastName, userName, Email, Password, AvatarURL string) (*User, error) {
	if firstName == "" || lastName == "" || userName == "" || Email == "" || Password == "" {
		return nil, errors.NewBadRequest("empty fields")
	}
	if len(firstName) > 50 || len(lastName) > 50 || len(userName) > 50 || len(Email) > 100 || len(Password) > 255 {
		return nil, errors.NewBadRequest("too long fields")
	}
	user := &User{FirstName: firstName, LastName: lastName, UserName: userName, Email: Email, PasswordHash: Password, UserStatus: Offline, AvatarURL: AvatarURL}
	return user, nil
}
