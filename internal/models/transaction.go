package models

import (
	"github.com/jinzhu/gorm"
)

type TransactionModel struct {
	gorm.Model
	AccountId uint
	Account   AccountModel `gorm:"foreignKey:AccountId"`
	Amount    float64
}

type Transaction struct {
	ID        uint    `json:"id,omitempty"`
	AccountId uint    `json:"account_id"`
	Amount    float64 `json:"amount" validate:"required,gte=0"`
	Status    string  `json:"status_msg,omitempty"`
}
