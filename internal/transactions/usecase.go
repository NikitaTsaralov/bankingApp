package transactions

import "github.com/NikitaTsaralov/bankingApp/internal/models"

type UseCase interface {
	PutMoney(trasaction *models.ResponseTransaction) (*models.ResponseTransaction, error)
	GetMoney(trasaction *models.ResponseTransaction) (*models.ResponseTransaction, error)
}
