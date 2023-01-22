package repository

import (
	"fmt"

	"github.com/NikitaTsaralov/bankingApp/internal/models"
	"github.com/NikitaTsaralov/bankingApp/internal/users"
	"github.com/jinzhu/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func Init(db *gorm.DB) users.Repository {
	return &userRepo{
		db: db,
	}
}

func (users *userRepo) Register(user *models.ResponseUser) (*models.ResponseUser, error) {
	tx := users.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, fmt.Errorf("complex transaction error: %v", err)
	}

	gormUser := &models.User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}
	if dbc := tx.Create(&gormUser); dbc.Error != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error gorm.DB.Create: %v", dbc.Error)
	}

	gormAccount := &models.Account{
		Type:    "Daily Account",
		Name:    string(user.Username + "'s" + " account"),
		Balance: 0,
		UserID:  gormUser.ID,
	}
	if dbc := tx.Create(&gormAccount); dbc.Error != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error creating account: %v", dbc.Error)
	}
	if dbc := tx.Commit(); dbc.Error != nil {
		return nil, fmt.Errorf("error creating: %v", dbc.Error)
	}

	return &models.ResponseUser{
		ID:       gormUser.ID,
		Username: gormUser.Username,
		Email:    gormUser.Email,
		Password: "********",
		Account: models.ResponseAccount{
			ID:      gormAccount.ID,
			Name:    gormAccount.Name,
			Balance: gormAccount.Balance,
			UserID:  gormAccount.UserID,
		},
	}, nil
}

func (users *userRepo) GetUserByName(username string) (*models.ResponseUser, error) {
	gormUser := &models.User{}
	if users.db.Table("users").Select("*").Where("users.username = ? ", username).First(&gormUser).RecordNotFound() {
		return nil, fmt.Errorf("user not found, try /register or check credentials")
	}

	gormAccount := &models.Account{}
	if users.db.Table("accounts").Select("*").Where("accounts.user_id = ? ", gormUser.ID).First(&gormAccount).RecordNotFound() {
		return nil, fmt.Errorf("account not found, contact DB administrator")
	}

	return &models.ResponseUser{
		ID:       gormUser.ID,
		Username: gormUser.Username,
		Email:    gormUser.Email,
		Password: gormUser.Password,
		Account: models.ResponseAccount{
			ID:      gormAccount.ID,
			Name:    gormAccount.Name,
			Balance: gormAccount.Balance,
			UserID:  gormAccount.UserID,
		},
	}, nil
}

func (users *userRepo) GetUserById(userId uint) (*models.ResponseUser, error) {
	gormUser := &models.User{}
	if users.db.Table("users").Select("*").Where("users.id = ? ", userId).First(&gormUser).RecordNotFound() {
		return nil, fmt.Errorf("user not found, try /register or check credentials")
	}

	gormAccount := &models.Account{}
	if users.db.Table("accounts").Select("*").Where("accounts.user_id = ? ", userId).First(&gormAccount).RecordNotFound() {
		return nil, fmt.Errorf("account not found, contact DB administrator")
	}

	return &models.ResponseUser{
		ID:       gormUser.ID,
		Username: gormUser.Username,
		Email:    gormUser.Email,
		Password: gormUser.Password,
		Account: models.ResponseAccount{
			ID:      gormAccount.ID,
			Name:    gormAccount.Name,
			Balance: gormAccount.Balance,
			UserID:  gormAccount.UserID,
		},
	}, nil
}

func (users *userRepo) GetAccountByUserId(userId uint) (*models.ResponseAccount, error) {
	gormAccount := &models.Account{}
	if users.db.Table("accounts").Select("*").Where("accounts.user_id = ? ", userId).First(&gormAccount).RecordNotFound() {
		return nil, fmt.Errorf("user not found, try /register or check credentials")
	}

	return &models.ResponseAccount{
		ID:      gormAccount.ID,
		Name:    gormAccount.Name,
		Balance: gormAccount.Balance,
		UserID:  gormAccount.UserID,
	}, nil
}

func (users *userRepo) GetTransactionsByUserId(userId uint) ([]models.ResponseTransaction, error) {
	gormTransaction := &models.Transaction{}
	gormTransactions := []models.Transaction{}

	users.db.Model(gormTransaction).Select("*").Joins("left join accounts on accounts.id = transactions.account_id").Joins("left join users on accounts.user_id = users.id").Where("accounts.user_id = ? ", userId).Scan(&gormTransactions)

	transactions := []models.ResponseTransaction{}
	for _, transaction := range gormTransactions {
		transactions = append(transactions, models.ResponseTransaction{
			ID:        transaction.ID,
			AccountId: transaction.AccountId,
			Amount:    transaction.Amount,
		})
	}
	return transactions, nil
}

func (users *userRepo) GetTransaction(id uint, userId uint) (*models.ResponseTransaction, error) {
	gormTransaction := &models.Transaction{}

	if users.db.Model(gormTransaction).Select("*").
		Joins("left join accounts on accounts.id = transactions.account_id").
		Joins("left join users on accounts.user_id = users.id").
		Where("accounts.user_id = ? ", userId).Where("transactions.id = ? ", id).
		Scan(&gormTransaction).RecordNotFound() {
		return nil, fmt.Errorf("not found")
	}

	return &models.ResponseTransaction{
		ID:        gormTransaction.ID,
		AccountId: gormTransaction.AccountId,
		Amount:    gormTransaction.Amount,
	}, nil
}
