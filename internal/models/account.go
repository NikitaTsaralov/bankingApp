package models

import "github.com/jinzhu/gorm"

type AccountModel struct {
	gorm.Model
	Type    string
	Name    string
	Balance float64
	UserID  uint
	User    UserModel `gorm:"foreignKey:UserID"`
}

type Account struct {
	ID      uint    `json:"id,omitempty" validate:"omitempty"`
	Balance float64 `json:"balance" validate:"required,gte=0"`
	UserID  uint    `json:"user_id" validate:"omitempty,required"`
}
