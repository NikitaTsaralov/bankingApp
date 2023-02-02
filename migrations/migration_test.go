package migrations

import (
	"testing"

	"github.com/NikitaTsaralov/bankingApp/config"
	"github.com/NikitaTsaralov/bankingApp/pkg/db"
)

func TestMigration(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _ := db.NewPsqlGormDB(tt.config)
			migration := &Migration{
				database: db,
			}
			err := migration.Migrate()
			if err != nil {
				t.Error("got error ", err)
			}
		})
	}
}
