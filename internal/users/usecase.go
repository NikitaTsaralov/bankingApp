package users

import (
	"github.com/NikitaTsaralov/bankingApp/internal/models"
)

type UseCase interface {
	Register(user *models.User) (*models.UserWithToken, error)
	Login(user *models.User) (*models.UserWithToken, error)

	GetById(id uint) (*models.User, error)
}
