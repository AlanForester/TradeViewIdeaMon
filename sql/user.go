package sql

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Uid       string
	Username  string
	Status    string
	IsBlocked bool
}

func (u User) Recipient() string {
	return u.Uid
}
