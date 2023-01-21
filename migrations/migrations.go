package migrations

import (
	"fmt"

	"github.com/NikitaTsaralov/bankingApp/config"
	"github.com/NikitaTsaralov/bankingApp/internal/models"
	"github.com/NikitaTsaralov/bankingApp/pkg/db"
	"github.com/jinzhu/gorm"
)

type Migration struct {
	database *gorm.DB
}

func (migration *Migration) Init(config *config.Config) (err error) {
	migration.database, err = db.Init(config)
	if err != nil {
		return fmt.Errorf("error init PsqlDB: %v", err)
	}
	return
}

func (migration *Migration) setup() error {
	transactionTypes := &[2]models.TransactionType{
		{
			Name: "AddSum",
		},
		{
			Name: "WithdrawSum",
		},
	}

	for _, transactionType := range transactionTypes {
		if dbc := migration.database.Create(&transactionType); dbc.Error != nil {
			return fmt.Errorf("gorm.DB.Create %v failed: %v", transactionType, dbc.Error)
		}
	}
	return nil
}

func (migration *Migration) Migrate() error {
	users := &models.User{}
	accounts := &models.Account{}
	transactionTypes := &models.TransactionType{}
	transactions := &models.Transaction{}

	// create tables
	if dbc := migration.database.AutoMigrate(&users, &accounts, &transactionTypes, &transactions); dbc.Error != nil {
		return fmt.Errorf("gorm.DB.AutoMigrate failed: %v", dbc.Error)
	}

	// fill tables
	if err := migration.setup(); err != nil {
		return fmt.Errorf("setup failed: %v", err)
	}
	return nil
}
