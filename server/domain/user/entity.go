package user

import (
	"fmt"
	"time"
)

type UserStatus int // симуляция enum

const (
	Offline UserStatus = iota // Отчет начинается с 0
	Online
)

type User struct {
	ID              int    // string or int зависит от того будет ли генерировать его сервер или бд
	FirstName       string // длина не больше 50
	LastName        string // длина не больше 50
	UserName        string // Буквы английского алфавита и _, должен быть уникальным в бд длина не больше 50
	Email           string // Длина не больше 100
	PasswordHash    string //длина не больше 255
	IsEmailVerified bool
	CreatedAt       time.Time
	AvatarURL       string
	LastSeenAt      time.Time
	UserStatus      UserStatus
}

func NewUser(firstName, lastName, userName, Email, Password, AvatarURL string) (*User, error) {
	if firstName == "" || lastName == "" || userName == "" || Email == "" || Password == "" {
		return nil, fmt.Errorf(`nil parameter`)
	}
	if len(firstName) > 50 || len(lastName) > 50 || len(userName) > 50 || len(Email) > 100 || len(Password) > 255 {
		return nil, fmt.Errorf(`too big parameter`)
	}
	user := &User{FirstName: firstName, LastName: lastName, UserName: userName, Email: Email, PasswordHash: Password, IsEmailVerified: false, UserStatus: 1, AvatarURL: AvatarURL}
	return user, nil
}
