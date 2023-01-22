package transactions

import "github.com/NikitaTsaralov/bankingApp/internal/models"

type UseCase interface {
	PublishMoneyOperation(userId uint, trasaction *models.ResponseTransaction) (*models.ResponseTransaction, error)
	MoneyOperation(trasaction *models.ResponseTransaction) (*models.ResponseTransaction, error)
}
