package sql

import "gorm.io/gorm"

type Author struct {
	gorm.Model
	Name      string `gorm:"uniqueIndex"`
	IsBlocked bool
}
