package middleware

import (
	"log"

	"github.com/NikitaTsaralov/bankingApp/config"
	"github.com/NikitaTsaralov/bankingApp/internal/users"
)

type MiddlewareManager struct {
	cfg    *config.Config
	userUC users.UseCase
	logger *log.Logger
}

func Init(cfg *config.Config, userUC users.UseCase, logger *log.Logger) *MiddlewareManager {
	return &MiddlewareManager{
		cfg:    cfg,
		userUC: userUC,
		logger: logger,
	}
}
