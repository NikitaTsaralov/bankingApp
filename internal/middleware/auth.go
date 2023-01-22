package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/NikitaTsaralov/bankingApp/pkg/httpErrors"
	"github.com/NikitaTsaralov/bankingApp/pkg/token"
	"github.com/labstack/echo/v4"
)

func (mw *MiddlewareManager) validateJWTToken(tokenString string, c echo.Context) error {
	userId, err := token.ValidateToken(tokenString, mw.cfg)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, httpErrors.InvalidJWTToken)
	}

	user, err := mw.authUC.GetUserById(userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, httpErrors.InvalidJWTToken)
	}

	c.Set("user", user.ID)
	return nil
}

func (mw *MiddlewareManager) AuthJWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			bearerHeader := c.Request().Header.Get("Authorization")
			if bearerHeader == "" {
				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
			}

			headerParts := strings.Split(bearerHeader, " ")
			if len(headerParts) != 2 {
				log.Println("auth middleware", fmt.Sprint("headerParts", "len(headerParts) != 2"))
				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
			}

			tokenString := headerParts[1]
			if err := mw.validateJWTToken(tokenString, c); err != nil {
				return c.JSON(http.StatusUnauthorized, err)
			}

			return next(c)
		}
	}
}
