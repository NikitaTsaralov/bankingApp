package models

import "github.com/jinzhu/gorm"

type Account struct {
	gorm.Model
	Type    string
	Name    string
	Balance float64
	UserID  uint
	User    User `gorm:"foreignKey:UserID"`
}

type ResponseAccount struct {
	ID      uint    `json:"id,omitempty"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
	UserID  uint    `json:"user_id"`
}
