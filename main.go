package main

import (
	"fmt"
	"log"
	"os"

	"github.com/NikitaTsaralov/bankingApp/config"
	"github.com/NikitaTsaralov/bankingApp/internal/server"
	"github.com/NikitaTsaralov/bankingApp/pkg/db"
	"github.com/NikitaTsaralov/bankingApp/pkg/logger"
	"github.com/NikitaTsaralov/bankingApp/pkg/rabbitmq"
	"github.com/NikitaTsaralov/bankingApp/pkg/utils"
)

func main() {
	fmt.Println("Starting REST API")

	// load and parse config file
	// TODO: this
	configPath := utils.GetConfigPath("config", os.Getenv("config"))
	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig cfgPath: %s failed: %v", configPath, err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig failed: %v", err)
	}

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

	database, err := db.NewPsqlGormDB(cfg)
	if err != nil {
		log.Fatalf("DB Init failed: %v", err)
	}

	broker, err := rabbitmq.Init(cfg)
	if err != nil {
		log.Fatalf("RabbitMQ Init failed: %v", err)
	}

	s := server.Init(cfg, database, broker, &log.Logger{})
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
