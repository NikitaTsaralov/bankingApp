package users

import (
	"github.com/NikitaTsaralov/bankingApp/internal/models"
)

type Repository interface {
	Register(user *models.ResponseUser) (*models.ResponseUser, error)
	GetUserByName(username string) (*models.ResponseUser, error)
	GetUserById(id uint) (*models.ResponseUser, error)

	GetAccountByUserId(id uint) (*models.ResponseAccount, error)
	GetTransactionsByUserId(id uint) ([]models.ResponseTransaction, error)
}
