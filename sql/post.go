package sql

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title  string
	Author string
	Tp     string
	Pair   string
	Date   string
	Image  string
	Video  string
	Descr  string
	Url    string `gorm:"uniqueIndex"`
}
