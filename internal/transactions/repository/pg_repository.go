package repository

import (
	"fmt"

	"github.com/NikitaTsaralov/bankingApp/internal/models"
	"github.com/NikitaTsaralov/bankingApp/internal/transactions"
	"github.com/jinzhu/gorm"
)

type transactionRepo struct {
	db *gorm.DB
}

func Init(db *gorm.DB) transactions.Repository {
	return &transactionRepo{
		db: db,
	}
}

func (transactions *transactionRepo) MoneyOperation(transaction *models.ResponseTransaction) (*models.ResponseTransaction, error) {
	tx := transactions.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, fmt.Errorf("complex transaction error: %v", err)
	}

	gormTransaction := &models.Transaction{
		AccountId: transaction.AccountId,
		Amount:    transaction.Amount,
	}
	if dbc := tx.Create(&gormTransaction); dbc.Error != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error creating account: %v", dbc.Error)
	}

	gormAccount := &models.Account{}
	if transactions.db.Table("accounts").Select("*").Where("accounts.id = ? ", transaction.AccountId).First(&gormAccount).RecordNotFound() {
		tx.Rollback()
		return nil, fmt.Errorf("account not found, try /register or check credentials")
	}

	gormAccount.Balance += transaction.Amount
	if gormAccount.Balance <= 0 {
		tx.Rollback()
		return &models.ResponseTransaction{
			ID:        gormTransaction.ID,
			AccountId: gormAccount.ID,
			Amount:    gormTransaction.Amount,
			Status:    "you're out of money, canceled",
		}, nil
	}

	if dbc := transactions.db.Save(&gormAccount); dbc.Error != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error updating Account: %v", dbc.Error)
	}

	if dbc := tx.Commit(); dbc.Error != nil {
		return nil, fmt.Errorf("error creating transaction: %v", dbc.Error)
	}

	return &models.ResponseTransaction{
		ID:        gormTransaction.ID,
		AccountId: gormAccount.ID,
		Amount:    gormTransaction.Amount,
	}, nil
}
