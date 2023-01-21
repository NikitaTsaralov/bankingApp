package main

import (
	"fmt"
	"log"

	"github.com/NikitaTsaralov/bankingApp/config"
	"github.com/NikitaTsaralov/bankingApp/internal/server"
	"github.com/NikitaTsaralov/bankingApp/pkg/db"
	"github.com/NikitaTsaralov/bankingApp/pkg/rabbitmq"
	"github.com/NikitaTsaralov/bankingApp/pkg/utils"
)

func main() {
	fmt.Println("Starting REST API")

	// load and parse config file
	// TODO: this
	// configPath := utils.GetConfigPath(os.Getenv("config"))
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

	broker := rabbitmq.Init(cfg, &log.Logger{})

	s := server.Init(cfg, database, broker, &log.Logger{})
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
