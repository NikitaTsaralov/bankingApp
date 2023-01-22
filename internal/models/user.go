package models

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

type UserWithToken struct {
	User  *ResponseUser `json:"user"`
	Token string        `json:"token"`
}

type ResponseUser struct {
	ID       uint            `json:"id,omitempty"`
	Username string          `json:"username" validate:"required"`
	Email    string          `json:"email" validate:"email,required"`
	Password string          `json:"password" validate:"required"`
	Account  ResponseAccount `json:"account,omitempty"`
}

func (user *ResponseUser) Validate() error {
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		return fmt.Errorf("Validation error: %v", err)
	}

	if len(user.Password) < 8 {
		return fmt.Errorf("Validation error: password too short")
	}
	return nil
}

func (user *ResponseUser) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}

func (user *ResponseUser) ComparePassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password)); err != nil {
		return err
	}
	return nil
}
