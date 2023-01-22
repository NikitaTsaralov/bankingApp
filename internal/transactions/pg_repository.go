package transactions

import "github.com/NikitaTsaralov/bankingApp/internal/models"

type Repository interface {
	MoneyOperation(transaction *models.ResponseTransaction) (*models.ResponseTransaction, error)
}
