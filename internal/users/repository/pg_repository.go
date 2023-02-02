package repository

import (
	"fmt"

	"github.com/NikitaTsaralov/bankingApp/internal/models"
	"github.com/NikitaTsaralov/bankingApp/internal/users"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type userRepo struct {
	db *gorm.DB
}

func Init(db *gorm.DB) users.Repository {
	return &userRepo{
		db: db,
	}
}

func (users *userRepo) Register(user *models.User) (*models.User, error) {
	tx := users.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, errors.Wrap(err, "userRepo.Register.TransactionRollback")
	}

	userModelImpl := &models.UserModel{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
	}
	if dbc := tx.Create(&userModelImpl); dbc.Error != nil {
		tx.Rollback()
		return nil, errors.Wrap(dbc.Error, "userRepo.Register.CreateUser")
	}

	accountModelImpl := &models.AccountModel{
		Balance: 0,
		UserID:  userModelImpl.ID,
	}
	if dbc := tx.Create(&accountModelImpl); dbc.Error != nil {
		tx.Rollback()
		return nil, errors.Wrap(dbc.Error, "userRepo.Register.CreateAccount")
	}
	if dbc := tx.Commit(); dbc.Error != nil {
		return nil, errors.Wrap(dbc.Error, "userRepo.Register.Commit")
	}

	return &models.User{
		ID:        userModelImpl.ID,
		FirstName: userModelImpl.FirstName,
		LastName:  userModelImpl.LastName,
		Email:     userModelImpl.Email,
		Password:  userModelImpl.Password,
	}, nil
}

func (users *userRepo) GetById(userId uint) (*models.User, error) {
	userModelImpl := &models.UserModel{}
	if users.db.Take(&userModelImpl).RecordNotFound() {
		return nil, fmt.Errorf("user not found, try /register or check credentials")
	}

	return &models.User{
		ID:        userModelImpl.ID,
		FirstName: userModelImpl.FirstName,
		LastName:  userModelImpl.LastName,
		Email:     userModelImpl.Email,
		Password:  userModelImpl.Password,
	}, nil
}

func (users *userRepo) GetByEmail(userEmail string) (*models.User, error) {
	userModelImpl := &models.UserModel{}
	if users.db.Where("email = ?", userEmail).Take(&userModelImpl).RecordNotFound() {
		return nil, fmt.Errorf("user not found, try /register or check credentials")
	}

	return &models.User{
		ID:        userModelImpl.ID,
		FirstName: userModelImpl.FirstName,
		LastName:  userModelImpl.LastName,
		Email:     userModelImpl.Email,
		Password:  userModelImpl.Password,
	}, nil
}
