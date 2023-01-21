package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/NikitaTsaralov/bankingApp/internal/models"
	"github.com/NikitaTsaralov/bankingApp/pkg/token"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func (s *Server) putMoney(c echo.Context) error {
	auth := c.Request().Header.Get("Authorization")
	userId, err := token.ValidateToken(auth)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("JWT Validation error: %v", err))
	}

	var transaction models.ResponseTransaction
	err = json.NewDecoder(c.Request().Body).Decode(&transaction)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Parse JSON error: %v", err))
	}

	validate := validator.New()
	err = validate.Struct(transaction)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Validation error: %v", err))
	}

	if transaction.Amount == 0.0 {
		return c.String(http.StatusBadRequest, "Validation error: cannot operate with 0 money amount")
	}

	account := &models.Account{}
	if s.db.Table("accounts").Select("*").Where("user_id = ? ", userId).First(&account).RecordNotFound() {
		return c.String(http.StatusNotFound, "Account not found")
	}

	transaction.AccountId = account.ID
	jsonBytes, err := json.Marshal(transaction)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("JSON Marshal error: %v", err))
	}

	err = s.broker.Send(jsonBytes)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Broker send failed: %v", err))
	}

	return c.JSON(http.StatusOK, &transaction)
}

func (s *Server) getMoney(c echo.Context) error {
	auth := c.Request().Header.Get("Authorization")
	userId, err := token.ValidateToken(auth)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("JWT Validation error: %v", err))
	}

	var transaction models.ResponseTransaction
	err = json.NewDecoder(c.Request().Body).Decode(&transaction)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Parse JSON error: %v", err))
	}

	validate := validator.New()
	err = validate.Struct(transaction)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Validation error: %v", err))
	}

	if transaction.Amount == 0.0 {
		return c.String(http.StatusBadRequest, "Validation error: cannot operate with 0 money amount")
	}

	account := &models.Account{}
	if s.db.Table("accounts").Select("*").Where("user_id = ? ", userId).First(&account).RecordNotFound() {
		return c.String(http.StatusNotFound, "Account not found for user, please contact administrator")
	}

	transaction.AccountId = account.ID
	transaction.Amount = -transaction.Amount
	jsonBytes, err := json.Marshal(transaction)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("JSON Marshal error: %v", err))
	}

	err = s.broker.Send(jsonBytes)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Broker send failed: %v", err))
	}

	return c.JSON(http.StatusOK, &transaction)
}
