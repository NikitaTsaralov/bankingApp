package main

import (
	"log"

	"github.com/NikitaTsaralov/bankingApp/config"
	"github.com/NikitaTsaralov/bankingApp/migrations"
	"github.com/NikitaTsaralov/bankingApp/pkg/utils"
)

func main() {
	configPath := utils.GetConfigPath("config", "local")
	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig cfgPath: %s failed: %v", configPath, err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig failed: %v", err)
	}

	migration, err := migrations.Init(cfg)
	if err != nil {
		log.Fatalf("Migration Init failed: %v", err)
	}

	err = migration.Migrate()
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}
