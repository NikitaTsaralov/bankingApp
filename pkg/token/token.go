package token

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/NikitaTsaralov/bankingApp/config"
	"github.com/NikitaTsaralov/bankingApp/internal/models"
	"github.com/dgrijalva/jwt-go"
)

func PrepareToken(user *models.ResponseUser, cfg *config.Config) (token string, err error) {
	tokenContent := jwt.MapClaims{
		"user_id":    strconv.Itoa(int(user.ID)),
		"account_id": strconv.Itoa(int(user.Account.ID)),
		"expiry":     time.Now().Add(time.Minute * time.Duration(cfg.Server.JwtExpire)).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err = jwtToken.SignedString([]byte(cfg.Server.JwtSecretKey))
	if err != nil {
		return "", fmt.Errorf("problem generating token: %v", err)
	}
	return token, nil
}

func ValidateToken(jwtToken string, cfg *config.Config) (userId uint, err error) {
	cleanJWT := strings.Replace(jwtToken, "Bearer ", "", -1)
	tokenData := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(cleanJWT, tokenData, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Server.JwtSecretKey), nil
	})
	if err != nil {
		return 0, fmt.Errorf("problem validating token: %v", err)
	}
	if token.Valid {
		if val, err := tokenData["user_id"].(string); err {
			uintVal, errParse := strconv.ParseUint(val, 10, 32)
			if errParse != nil {
				return 0, fmt.Errorf("problem parsing token data token: %v", err)
			}
			return uint(uintVal), nil
		} else {
			return 0, fmt.Errorf("token invalid: %v", err)
		}
	} else {
		return 0, fmt.Errorf("token invalid: %v", err)
	}
}
