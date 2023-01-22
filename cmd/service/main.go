package main

import (
	"fmt"
	"log"

	"github.com/NikitaTsaralov/bankingApp/config"
	"github.com/NikitaTsaralov/bankingApp/pkg/db"
	"github.com/NikitaTsaralov/bankingApp/pkg/utils"

	"github.com/NikitaTsaralov/bankingApp/internal/transactions/delivery/rabbitmq"
	transactionRepo "github.com/NikitaTsaralov/bankingApp/internal/transactions/repository"
	transactionUseCase "github.com/NikitaTsaralov/bankingApp/internal/transactions/usecase"
	userRepo "github.com/NikitaTsaralov/bankingApp/internal/users/repository"
	// amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting service")

	configPath := utils.GetConfigPath("config", "local")
	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig cfgPath: %s failed: %v", configPath, err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig failed: %v", err)
	}

	database, err := db.Init(cfg)
	if err != nil {
		log.Fatalf("DB Init failed: %v", err)
	}

	logger := &log.Logger{}
	uRepo := userRepo.Init(database)
	tRepo := transactionRepo.Init(database)
	publisher, err := rabbitmq.InitTransactionPublisher(cfg, logger)
	if err != nil {
		log.Fatalf("InitTransactionPublisher failed: %v", err)
	}
	transactionUseCase := transactionUseCase.Init(cfg, tRepo, uRepo, publisher, logger)
	consumer, err := rabbitmq.InitTransactionConsumer(cfg, transactionUseCase, logger)
	if err != nil {
		log.Fatalf("InitTransactionConsumer failed: %v", err)
	}

	err = consumer.StartConsumer()
	if err != nil {
		log.Fatalf("Worker run failed: %v", err)
	}
}
