package http

import (
	"github.com/NikitaTsaralov/bankingApp/internal/middleware"
	"github.com/NikitaTsaralov/bankingApp/internal/transactions"
	"github.com/labstack/echo/v4"
)

func MapTransactionRoutes(usersGroup *echo.Group, h transactions.Handlers, mw *middleware.MiddlewareManager) {
	usersGroup.Use(mw.AuthJWTMiddleware())
	usersGroup.POST("/get_money", h.GetMoney())
	usersGroup.POST("/put_money", h.PutMoney())
}
