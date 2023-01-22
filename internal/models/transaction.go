package models

import (
	"fmt"

	"github.com/go-playground/validator/v10"
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

func (transaction *ResponseTransaction) Validate() error {
	validate := validator.New()
	err := validate.Struct(transaction)
	if err != nil {
		return fmt.Errorf("validation error: %v", err)
	}

	if transaction.Amount <= 0 {
		return fmt.Errorf("validation error: you cannot put %v money", transaction.Amount)
	}

	return nil
}
