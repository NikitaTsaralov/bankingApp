package transactions

import "github.com/NikitaTsaralov/bankingApp/internal/models"

type UseCase interface {
	MoneyOperation(userId uint, trasaction *models.ResponseTransaction) (*models.ResponseTransaction, error)
}
