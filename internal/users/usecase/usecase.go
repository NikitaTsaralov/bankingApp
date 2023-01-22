package usecase

import (
	"fmt"
	"log"

	"github.com/NikitaTsaralov/bankingApp/config"
	"github.com/NikitaTsaralov/bankingApp/internal/models"
	"github.com/NikitaTsaralov/bankingApp/internal/users"
	"github.com/NikitaTsaralov/bankingApp/pkg/token"
)

type userUC struct {
	cfg      *config.Config
	userRepo users.Repository
	logger   *log.Logger
}

func Init(cfg *config.Config, usersRepo users.Repository, logger *log.Logger) users.UseCase {
	return &userUC{
		cfg:      cfg,
		userRepo: usersRepo,
		logger:   logger,
	}
}

func (uc *userUC) Register(user *models.ResponseUser) (*models.UserWithToken, error) {
	user, err := uc.userRepo.Register(user)
	if err != nil {
		return nil, fmt.Errorf("error transactionsRepo.Register: %v", err)
	}

	token, err := token.PrepareToken(user, uc.cfg)
	if err != nil {
		return nil, fmt.Errorf("error token.PrepareToken: %v", err)
	}

	return &models.UserWithToken{
		User:  user,
		Token: token,
	}, nil
}

func (uc *userUC) Login(user *models.ResponseUser) (*models.UserWithToken, error) {
	repoUser, err := uc.userRepo.GetUserByName(user.Username)
	if err != nil {
		return nil, fmt.Errorf("error transactionsRepo.Login: %v", err)
	}

	if err := user.ComparePassword(repoUser.Password); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	user.Password = "********"

	token, err := token.PrepareToken(repoUser, uc.cfg)
	if err != nil {
		return nil, fmt.Errorf("error token.PrepareToken: %v", err)
	}

	return &models.UserWithToken{
		User:  repoUser,
		Token: token,
	}, nil
}

func (uc *userUC) GetUserById(id uint) (*models.ResponseUser, error) {
	repoUser, err := uc.userRepo.GetUserById(id)
	if err != nil {
		return nil, fmt.Errorf("error userRepo.GetById: %v", err)
	}
	return repoUser, nil
}

func (uc *userUC) GetAccountByUserId(userId uint) (*models.ResponseAccount, error) {
	repoAccounts, err := uc.userRepo.GetAccountByUserId(userId)
	if err != nil {
		return nil, fmt.Errorf("error userRepo.GetAccountByUserId: %v", err)
	}
	return repoAccounts, nil
}

func (uc *userUC) GetTransactionsByUserId(userId uint) ([]models.ResponseTransaction, error) {
	repoTransactions, err := uc.userRepo.GetTransactionsByUserId(userId)
	if err != nil {
		return nil, fmt.Errorf("error userRepo.GetTransactionsByUserId: %v", err)
	}
	return repoTransactions, nil
}

func (uc *userUC) GetTransaction(id uint, userId uint) (*models.ResponseTransaction, error) {
	repoTransaction, err := uc.userRepo.GetTransaction(id, userId)
	if err != nil {
		return nil, fmt.Errorf("error userRepo.GetTransactionsByUserId: %v", err)
	}
	return repoTransaction, nil
}
