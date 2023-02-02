package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/NikitaTsaralov/bankingApp/pkg/httpErrors"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

func (mw *MiddlewareManager) validateJWTToken(tokenString string, c echo.Context) error {
	if tokenString == "" {
		return httpErrors.InvalidJWTToken
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method %v", token.Header["alg"])
		}
		secret := []byte(mw.cfg.Server.JwtSecretKey)
		return secret, nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return httpErrors.InvalidJWTToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["id"].(string)
		if !ok {
			return httpErrors.InvalidJWTClaims
		}

		userUintID, err := strconv.ParseUint(userID, 10, 32)
		if err != nil {
			return err
		}

		u, err := mw.userUC.GetById(uint(userUintID))
		if err != nil {
			return err
		}

		c.Set("user", u)
	}
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
