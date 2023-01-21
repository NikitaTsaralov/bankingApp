package models

import "github.com/jinzhu/gorm"

type Account struct {
	gorm.Model
	Type    string
	Name    string
	Balance uint
	UserID  uint
}
