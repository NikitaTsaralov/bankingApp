package transactions

import "github.com/NikitaTsaralov/bankingApp/internal/models"

type TransactionPublisher interface {
	Publish(*models.ResponseTransaction) (*models.ResponseTransaction, error)
}

type TransactionConsumer interface {
	StartConsumer(workerPoolSize int, queueName string) error
}
