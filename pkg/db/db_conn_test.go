package db

import (
	"testing"

	"github.com/NikitaTsaralov/bankingApp/config"
)

// Test only to enshure // local
func TestNewPsqlDB(t *testing.T) {
	tests := []struct {
		name   string
		config *config.Config
		gotErr bool
	}{
		{
			name: "success",
			config: &config.Config{
				Postgres: config.PostgresConfig{
					PostgresqlHost:     "127.0.0.1",
					PostgresqlPort:     "5432",
					PostgresqlUser:     "postgres",
					PostgresqlPassword: "postgres",
					PostgresqlDbname:   "auth_db",
					PostgresqlSSLMode:  false,
					PgDriver:           "pgx",
				},
			},
			gotErr: false,
		},
		{
			name: "error",
			config: &config.Config{
				Postgres: config.PostgresConfig{
					PostgresqlHost:     "127.0.0.1",
					PostgresqlPort:     "5433",
					PostgresqlUser:     "postgres",
					PostgresqlPassword: "postgres",
					PostgresqlDbname:   "auth_db",
					PostgresqlSSLMode:  false,
					PgDriver:           "pgx",
				},
			},
			gotErr: true,
		},
		{
			name: "error",
			config: &config.Config{
				Postgres: config.PostgresConfig{
					PostgresqlHost:     "192.168.10.135",
					PostgresqlPort:     "5432",
					PostgresqlUser:     "postgres",
					PostgresqlPassword: "postgres",
					PostgresqlDbname:   "auth_db",
					PostgresqlSSLMode:  false,
					PgDriver:           "pgx",
				},
			},
			gotErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Init(tt.config)
			if err != nil != tt.gotErr {
				t.Error(
					"expected gotErr: ", tt.gotErr, " got: ", err != nil,
				)
			}
		})
	}
}
