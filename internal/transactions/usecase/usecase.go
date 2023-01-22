package usecase

import (
	"log"

	"github.com/NikitaTsaralov/bankingApp/config"
	"github.com/NikitaTsaralov/bankingApp/internal/models"
	"github.com/NikitaTsaralov/bankingApp/internal/transactions"
)

type transactionsUC struct {
	cfg             *config.Config
	transactionRepo transactions.Repository
	logger          *log.Logger
}

func Init(cfg *config.Config, transactionRepo transactions.Repository, logger *log.Logger) *transactionsUC {
	return &transactionsUC{
		cfg:             cfg,
		transactionRepo: transactionRepo,
		logger:          logger,
	}
}

func (transactions *transactionsUC) PutMoney(trasaction *models.ResponseTransaction) (*models.ResponseTransaction, error) {
	return &models.ResponseTransaction{}, nil
}

func (transactions *transactionsUC) GetMoney(trasaction *models.ResponseTransaction) (*models.ResponseTransaction, error) {
	return &models.ResponseTransaction{}, nil
}
