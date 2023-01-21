package users

import (
	"github.com/NikitaTsaralov/bankingApp/internal/models"
	"github.com/NikitaTsaralov/bankingApp/pkg/token"
)

// Refactor prepareResponse
func PrepareResponse(user *models.User, account models.Account, withToken bool) (response map[string]interface{}, err error) {
	responseUser := &models.ResponseUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Password: "********",
		Account: models.ResponseAccount{
			ID:      account.ID,
			Name:    account.Name,
			Balance: account.Balance,
			UserID:  account.UserID,
		},
	}
	response = map[string]interface{}{}
	if withToken {
		token, err := token.PrepareToken(responseUser)
		if err != nil {
			return nil, err
		}
		response["jwt"] = token
	}
	response["data"] = responseUser
	return response, nil
}
