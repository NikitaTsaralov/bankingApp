package usecase

import (
	"log"
	"net/http"

	"github.com/NikitaTsaralov/bankingApp/config"
	"github.com/NikitaTsaralov/bankingApp/internal/models"
	"github.com/NikitaTsaralov/bankingApp/internal/users"
	"github.com/NikitaTsaralov/bankingApp/pkg/httpErrors"
	"github.com/NikitaTsaralov/bankingApp/pkg/utils"
	"github.com/pkg/errors"
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

func (uc *userUC) Register(user *models.User) (*models.UserWithToken, error) {
	err := utils.ValidateStruct(user)
	if err != nil {
		return nil, err
	}

	existsUser, err := uc.userRepo.GetByEmail(user.Email)
	if existsUser != nil || err == nil {
		return nil, httpErrors.NewRestErrorWithMessage(http.StatusBadRequest, httpErrors.ErrEmailAlreadyExists, nil)
	}

	err = user.PrepareCreate()
	if err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "userUC.Register.PrepareCreate"))
	}

	createdUser, err := uc.userRepo.Register(user)
	if err != nil {
		return nil, err
	}
	createdUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(createdUser, uc.cfg)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "userUC.Register.GenerateJWTToken"))
	}

	return &models.UserWithToken{
		User:  createdUser,
		Token: token,
	}, nil
}

func (uc *userUC) Login(user *models.User) (*models.UserWithToken, error) {
	err := utils.ValidateStruct(user)
	if err != nil {
		return nil, err
	}

	foundUser, err := uc.userRepo.GetByEmail(user.Email)
	if err != nil {
		return nil, errors.Wrap(err, "userUC.Login.GetByEmail")
	}

	if err := user.ComparePassword(foundUser.Password); err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.Wrap(err, "userUC.Login.ComparePassword"))
	}
	foundUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(foundUser, uc.cfg)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "userUC.Register.GenerateJWTToken"))
	}

	return &models.UserWithToken{
		User:  foundUser,
		Token: token,
	}, nil
}

func (uc *userUC) GetById(id uint) (*models.User, error) {
	foundUser, err := uc.userRepo.GetById(id)
	if err != nil {
		return nil, errors.Wrap(err, "userUC.Login.GetById")
	}
	foundUser.SanitizePassword()

	return foundUser, nil
}
