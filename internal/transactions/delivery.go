package transactions

import "github.com/labstack/echo/v4"

type Handlers interface {
	PutMoney() echo.HandlerFunc
	GetMoney() echo.HandlerFunc
}
