package usecase

import (
	"fmt"
	"log"

	"github.com/NikitaTsaralov/bankingApp/config"
	"github.com/NikitaTsaralov/bankingApp/internal/models"
	"github.com/NikitaTsaralov/bankingApp/internal/transactions"
	"github.com/NikitaTsaralov/bankingApp/internal/users"
)

type transactionsUC struct {
	cfg             *config.Config
	transactionRepo transactions.Repository
	userRepo        users.Repository
	publisher       transactions.TransactionPublisher
	logger          *log.Logger
}

func Init(cfg *config.Config, transactionRepo transactions.Repository, userRepo users.Repository, publisher transactions.TransactionPublisher, logger *log.Logger) *transactionsUC {
	return &transactionsUC{
		cfg:             cfg,
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
		publisher:       publisher,
		logger:          logger,
	}
}

func (transactions *transactionsUC) PublishMoneyOperation(userId uint, transaction *models.ResponseTransaction) (*models.ResponseTransaction, error) {
	// get account Id
	account, err := transactions.userRepo.GetAccountByUserId(userId)
	if err != nil {
		return nil, fmt.Errorf("error userRepo.GetAccountByUserId: %v", err)
	}

	// create msg for broker
	transaction.AccountId = account.ID
	resp, err := transactions.publisher.Publish(transaction)
	if err != nil {
		return nil, fmt.Errorf("broker send failed: %v", err)
	}

	return resp, nil
}

func (transactions *transactionsUC) MoneyOperation(transaction *models.ResponseTransaction) (*models.ResponseTransaction, error) {
	resp, err := transactions.transactionRepo.MoneyOperation(transaction)
	if err != nil {
		return nil, fmt.Errorf("error  transactionRepo.MoneyOperation: %v", err)
	}
	return resp, err
}
