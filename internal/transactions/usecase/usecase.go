package usecase

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/NikitaTsaralov/bankingApp/config"
	"github.com/NikitaTsaralov/bankingApp/internal/models"
	"github.com/NikitaTsaralov/bankingApp/internal/transactions"
	"github.com/NikitaTsaralov/bankingApp/internal/users"
	"github.com/NikitaTsaralov/bankingApp/pkg/rabbitmq"
)

type transactionsUC struct {
	cfg             *config.Config
	transactionRepo transactions.Repository
	userRepo        users.Repository
	broker          *rabbitmq.RabbitMQClient
	logger          *log.Logger
}

func Init(cfg *config.Config, transactionRepo transactions.Repository, userRepo users.Repository, broker *rabbitmq.RabbitMQClient, logger *log.Logger) *transactionsUC {
	return &transactionsUC{
		cfg:             cfg,
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
		broker:          broker,
		logger:          logger,
	}
}

func (transactions *transactionsUC) MoneyOperation(userId uint, transaction *models.ResponseTransaction) (*models.ResponseTransaction, error) {
	// get account Id
	account, err := transactions.userRepo.GetAccountByUserId(userId)
	if err != nil {
		return nil, fmt.Errorf("error userRepo.GetAccountByUserId: %v", err)
	}

	// create msg for broker
	transaction.AccountId = account.ID

	jsonBytes, err := json.Marshal(transaction)
	if err != nil {
		return nil, fmt.Errorf("JSON Marshal error: %v", err)
	}

	err = transactions.broker.Send(jsonBytes)
	if err != nil {
		return nil, fmt.Errorf("broker send failed: %v", err)
	}

	return &models.ResponseTransaction{}, nil
}
