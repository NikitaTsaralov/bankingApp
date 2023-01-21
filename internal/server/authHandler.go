package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/NikitaTsaralov/bankingApp/internal/models"
	"github.com/NikitaTsaralov/bankingApp/internal/users"
	"github.com/labstack/echo/v4"
)

func (auth *Server) Login(c echo.Context) error {
	var user models.ResponseUser
	err := json.NewDecoder(c.Request().Body).Decode(&user)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Parse JSON error: %v", err.Error()))
	}

	if err := user.Validate(); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	gormUser := &models.User{}
	if auth.db.Where("username = ? ", user.Username).First(&gormUser).RecordNotFound() {
		return c.String(http.StatusNotFound, "User not found, try /register or check credentials")
	}

	if err := user.ComparePassword(gormUser.Password); err != nil {
		return c.String(http.StatusForbidden, "Wrong password")
	}

	account := &models.Account{}
	if auth.db.Table("accounts").Select("*").Where("user_id = ? ", gormUser.ID).First(&account).RecordNotFound() {
		return c.String(http.StatusNotFound, "Account not found for user, please contact administrator")
	}

	response, err := users.PrepareResponse(gormUser, *account, true)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Error preparing response: %v", err.Error()))
	}

	return c.JSON(http.StatusOK, response)
}

func (auth *Server) Register(c echo.Context) error {
	var user models.ResponseUser
	err := json.NewDecoder(c.Request().Body).Decode(&user)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Parse JSON error: %v", err.Error()))
	}

	if err := user.Validate(); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	err = user.HashPassword()
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Error hashing password: %v", err.Error()))
	}

	gormUser := &models.User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}
	if dbc := auth.db.Create(&gormUser); dbc.Error != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error creating user: %v", dbc.Error))
	}

	account := &models.Account{
		Type:    "Daily Account",
		Name:    string(user.Username + "'s" + " account"),
		Balance: 0,
		UserID:  gormUser.ID,
	}
	if dbc := auth.db.Create(&account); dbc.Error != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error creating account: %v", dbc.Error))
	}

	response, err := users.PrepareResponse(gormUser, *account, true)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Error preparing response: %v", err.Error()))
	}

	return c.JSON(http.StatusOK, response)
}
