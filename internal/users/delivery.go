package users

import "github.com/labstack/echo/v4"

type Handlers interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc

	GetMe() echo.HandlerFunc
	GetMyAccount() echo.HandlerFunc
	History() echo.HandlerFunc
	GetTranaction() echo.HandlerFunc
}
