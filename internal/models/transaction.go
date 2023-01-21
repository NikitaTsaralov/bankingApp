package models

import (
	"github.com/jinzhu/gorm"
)

type Transaction struct {
	gorm.Model
	AccountId uint
	Account   Account `gorm:"foreignKey:AccountId"`
	Amount    float64
}

type ResponseTransaction struct {
	ID        uint    `json:"id,omitempty"`
	AccountId uint    `json:"account_id"`
	Amount    float64 `json:"amount" validate:"required"`
}
