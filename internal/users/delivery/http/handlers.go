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
		var user models.User
		err := json.NewDecoder(c.Request().Body).Decode(&user)
		if err != nil {
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, fmt.Sprintf("parse JSON error: %v", err), err))
		}

		createdUser, err := h.userUC.Register(&user)
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, createdUser)
	}
}

func (h *userHandlers) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var user models.User
		err := json.NewDecoder(c.Request().Body).Decode(&user)
		if err != nil {
			return c.JSON(http.StatusBadRequest, httpErrors.NewRestError(http.StatusBadRequest, fmt.Sprintf("parse JSON error: %v", err), err))
		}

		userWithToken, err := h.userUC.Login(&user)
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, userWithToken)
	}
}

func (h *userHandlers) GetMe() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(*models.User)
		if !ok {
			return c.JSON(httpErrors.ErrorResponse(httpErrors.NewUnauthorizedError(httpErrors.Unauthorized)))
		}

		return c.JSON(http.StatusOK, user)
	}
}
