package models

import (
	"strings"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	gorm.Model
	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
}

type User struct {
	ID        uint   `json:"id,omitempty" validate:"omitempty"`
	FirstName string `json:"first_name" validate:"lte=30"`
	LastName  string `json:"last_name" validate:"lte=30"`
	Email     string `json:"email" validate:"required,lte=60,email"`
	Password  string `json:"password,omitempty" validate:"omitempty,required,gte=8"`
}

type UserWithToken struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

func (u *User) PrepareCreate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Password = strings.TrimSpace(u.Password)

	if err := u.HashPassword(); err != nil {
		return err
	}

	return nil
}

// Sanitize user password
func (u *User) SanitizePassword() {
	u.Password = ""
}

// Hash user password with bcrypt
func (user *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}

// Compare user password and payload
func (user *User) ComparePassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password)); err != nil {
		return err
	}
	return nil
}
