package db

import (
	"fmt"

	"github.com/NikitaTsaralov/bankingApp/config"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

const (
	maxOpenConns    = 60
	connMaxLifetime = 120
	maxIdleConns    = 30
	connMaxIdleTime = 20
)

// Create InitDatabase function
func NewPsqlGormDB(c *config.Config) (DB *gorm.DB, err error) {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		c.Postgres.PostgresqlHost,
		c.Postgres.PostgresqlPort,
		c.Postgres.PostgresqlUser,
		c.Postgres.PostgresqlDbname,
		c.Postgres.PostgresqlPassword,
	)

	database, err := gorm.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error gorm.Open host=%s port=%s user=%s dbname=%s: %v",
			c.Postgres.PostgresqlHost,
			c.Postgres.PostgresqlPort,
			c.Postgres.PostgresqlUser,
			c.Postgres.PostgresqlDbname,
			err,
		)
	}

	// Set up connection pool
	database.DB().SetMaxOpenConns(maxOpenConns)
	database.DB().SetConnMaxLifetime(connMaxLifetime)
	database.DB().SetMaxIdleConns(connMaxIdleTime)
	database.DB().SetConnMaxIdleTime(connMaxIdleTime)

	if err = database.DB().Ping(); err != nil {
		return nil, fmt.Errorf("error db ping: %v", err)
	}

	return database, nil
}
