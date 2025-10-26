package user

import (
	"time"
)

type UserStatus int // симуляция enum

const (
	Offline UserStatus = iota // Отчет начинается с 0
	Online
)

type User struct {
	ID              int // string or int зависит от того будет ли генерировать его сервер или бд
	FirstName       string
	LastName        string
	UserName        string // Буквы английского алфавита и _, должен быть уникальным
	Email           string
	PasswordHash    string
	IsEmailVerified bool
	CreatedAt       time.Time
	AvatarURL       string
	LastSeenAt      time.Time
	UserStatus      UserStatus
}
