package main

import (
	"fmt"
	"log"

	"github.com/NikitaTsaralov/bankingApp/config"
	"github.com/NikitaTsaralov/bankingApp/pkg/db"
	"github.com/NikitaTsaralov/bankingApp/pkg/utils"
	"github.com/NikitaTsaralov/bankingApp/service"
	// amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting service")

	configPath := utils.GetConfigPath("config-service", "local")
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

	dumper := service.Init(cfg, database, &log.Logger{})

	err = dumper.Run()
	if err != nil {
		log.Fatalf("Worker run failed: %v", err)
	}
}
