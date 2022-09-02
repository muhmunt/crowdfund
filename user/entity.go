package user

import "time"

type User struct {
	ID             int
	Name           string
	Occupation     string
	Email          string
	PasswordHash   string
	AvatarFileName string
	Role           string
	Created_at     time.Time
	Updated_at     time.Time
}

// func GetUser()
