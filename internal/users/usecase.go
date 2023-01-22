package users

import (
	"github.com/NikitaTsaralov/bankingApp/internal/models"
)

type UseCase interface {
	Register(user *models.ResponseUser) (*models.UserWithToken, error)
	Login(user *models.ResponseUser) (*models.UserWithToken, error)
	GetUserById(id uint) (*models.ResponseUser, error)

	GetAccountByUserId(id uint) (*models.ResponseAccount, error)
	GetTransactionsByUserId(id uint) ([]models.ResponseTransaction, error)
}
