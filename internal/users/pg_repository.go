package users

import (
	"github.com/NikitaTsaralov/bankingApp/internal/models"
)

type Repository interface {
	Register(user *models.User) (*models.User, error)

	GetById(userID uint) (*models.User, error)
	GetByEmail(userEmail string) (*models.User, error)
}
