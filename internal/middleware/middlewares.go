package middleware

import (
	"log"

	"github.com/NikitaTsaralov/bankingApp/config"
	"github.com/NikitaTsaralov/bankingApp/internal/users"
)

type MiddlewareManager struct {
	cfg    *config.Config
	authUC users.UseCase
	logger *log.Logger
}

func Init(cfg *config.Config, authUC users.UseCase, logger *log.Logger) *MiddlewareManager {
	return &MiddlewareManager{
		cfg:    cfg,
		authUC: authUC,
		logger: logger,
	}
}
