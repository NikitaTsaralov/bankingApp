package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/NikitaTsaralov/bankingApp/config"
	"github.com/NikitaTsaralov/bankingApp/internal/models"
	"github.com/NikitaTsaralov/bankingApp/internal/users"
	"github.com/NikitaTsaralov/bankingApp/pkg/httpErrors"
	"github.com/labstack/echo/v4"
)

type userHandlers struct {
	cfg    *config.Config
	userUC users.UseCase
	logger *log.Logger
}

func Init(cfg *config.Config, userUC users.UseCase, logger *log.Logger) users.Handlers {
	return &userHandlers{
		cfg:    cfg,
		userUC: userUC,
		logger: logger,
	}
}

func (h *userHandlers) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		var user models.ResponseUser
		err := json.NewDecoder(c.Request().Body).Decode(&user)
		if err != nil {
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, fmt.Sprintf("parse JSON error: %v", err), err))
		}

		if err := user.Validate(); err != nil {
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, err.Error(), err))
		}

		err = user.HashPassword()
		if err != nil {
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, fmt.Sprintf("error hashing password: %v", err), err))
		}

		userWithToken, err := h.userUC.Register(&user)
		if err != nil {
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, fmt.Sprintf("error register: %v", err), err))
		}

		return c.JSON(http.StatusOK, userWithToken)
	}
}

func (h *userHandlers) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var user models.ResponseUser
		err := json.NewDecoder(c.Request().Body).Decode(&user)
		if err != nil {
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, fmt.Sprintf("parse JSON error: %v", err), err))
		}

		if err := user.Validate(); err != nil {
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, err.Error(), err))
		}

		userWithToken, err := h.userUC.Login(&user)
		if err != nil {
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, fmt.Sprintf("error login: %v", err), err))
		}

		return c.JSON(http.StatusOK, userWithToken)
	}
}

func (h *userHandlers) GetMe() echo.HandlerFunc {
	return func(c echo.Context) error {
		res := c.Get("user")

		if _, exists := res.(uint); !exists {
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, "wrong JWT token: `user` not in token", nil))
		}

		user, err := h.userUC.GetUserById(res.(uint))
		if err != nil {
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, fmt.Sprintf("error get user with id %d: %v", res.(uint), err), err))
		}

		return c.JSON(http.StatusOK, user)
	}
}

func (h *userHandlers) GetMyAccount() echo.HandlerFunc {
	return func(c echo.Context) error {
		res := c.Get("user")

		if _, exists := res.(uint); !exists {
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, "wrong JWT token: `user` not in token", nil))
		}

		account, err := h.userUC.GetAccountByUserId(res.(uint))
		if err != nil {
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, fmt.Sprintf("error getting account with user_id %d: %v", res.(uint), err), err))
		}

		return c.JSON(http.StatusOK, account)
	}
}

func (h *userHandlers) History() echo.HandlerFunc {
	return func(c echo.Context) error {
		res := c.Get("user")

		if _, exists := res.(uint); !exists {
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, "wrong JWT token: `user` not in token", nil))
		}

		transactions, err := h.userUC.GetTransactionsByUserId(res.(uint))
		if err != nil {
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, fmt.Sprintf("error getting transactions with user_id %d: %v", res.(uint), err), err))
		}

		return c.JSON(http.StatusOK, transactions)
	}
}
