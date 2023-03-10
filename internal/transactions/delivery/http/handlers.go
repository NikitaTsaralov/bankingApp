package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/NikitaTsaralov/bankingApp/config"
	"github.com/NikitaTsaralov/bankingApp/internal/models"
	"github.com/NikitaTsaralov/bankingApp/internal/transactions"
	"github.com/NikitaTsaralov/bankingApp/pkg/httpErrors"
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
		res := c.Get("user")
		if _, exists := res.(uint); !exists {
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, "wrong JWT token: `user` not in token", nil))
		}

		var transaction models.ResponseTransaction
		err := json.NewDecoder(c.Request().Body).Decode(&transaction)
		if err != nil {
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, fmt.Sprintf("parse JSON error: %v", err), err))
		}

		if err := transaction.Validate(); err != nil {
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, err.Error(), err))
		}

		transactionResp, err := transactions.transactionsUC.PublishMoneyOperation(res.(uint), &transaction)
		if err != nil {
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, fmt.Sprintf("error transactionsUC.PutMoney: %v", err), err))
		}

		return c.JSON(http.StatusOK, transactionResp)
	}
}

func (transactions *transactionHandlers) GetMoney() echo.HandlerFunc {
	return func(c echo.Context) error {
		res := c.Get("user")
		if _, exists := res.(uint); !exists {
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, "wrong JWT token: `user` not in token", nil))
		}

		var transaction models.ResponseTransaction
		err := json.NewDecoder(c.Request().Body).Decode(&transaction)
		if err != nil {
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, fmt.Sprintf("parse JSON error: %v", err), err))
		}

		if err := transaction.Validate(); err != nil {
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, err.Error(), err))
		}

		transaction.Amount = -transaction.Amount
		transactionResp, err := transactions.transactionsUC.PublishMoneyOperation(res.(uint), &transaction)
		if err != nil {
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, fmt.Sprintf("error transactionsUC.PutMoney: %v", err), err))
		}

		return c.JSON(http.StatusOK, transactionResp)
	}
}
