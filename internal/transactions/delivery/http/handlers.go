package http

import (
	"log"
	"net/http"

	"github.com/NikitaTsaralov/bankingApp/config"
	"github.com/NikitaTsaralov/bankingApp/internal/transactions"
	"github.com/labstack/echo/v4"
)

type transactionHandlers struct {
	cfg            *config.Config
	transactionsUC transactions.UseCase
	logger         *log.Logger
}

func Init(cfg *config.Config, transactionsUC transactions.UseCase, logger *log.Logger) *transactionHandlers {
	return &transactionHandlers{
		cfg:            cfg,
		transactionsUC: transactionsUC,
		logger:         logger,
	}
}

func (transactions *transactionHandlers) PutMoney() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"a": "b"})
	}
}

func (transactions *transactionHandlers) GetMoney() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"a": "b"})
	}
}
