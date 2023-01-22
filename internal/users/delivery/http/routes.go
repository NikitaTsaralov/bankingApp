package http

import (
	"github.com/NikitaTsaralov/bankingApp/internal/middleware"
	"github.com/NikitaTsaralov/bankingApp/internal/users"
	"github.com/labstack/echo/v4"
)

func MapAuthRoutes(usersGroup *echo.Group, h users.Handlers, mw *middleware.MiddlewareManager) {
	usersGroup.POST("/register", h.Register())
	usersGroup.POST("/login", h.Login())

	usersGroup.Use(mw.AuthJWTMiddleware())
	usersGroup.GET("/me", h.GetMe())
	usersGroup.GET("/my_account", h.GetMyAccount())

}

func MapHistoryRoutes(usersGroup *echo.Group, h users.Handlers, mw *middleware.MiddlewareManager) {
	usersGroup.Use(mw.AuthJWTMiddleware())
	usersGroup.GET("/operation_history", h.History())
	usersGroup.GET("/get_transaction", h.GetTranaction())
}
