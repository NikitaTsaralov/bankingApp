package config

type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	RabbitMQ RabbitConfig
}

type ServerConfig struct {
	AppVersion   string
	Port         string
	PprofPort    string
	Mode         string
	JwtSecretKey string
	CookieName   string
	SSL          bool
	CSRF         bool
	Debug        bool
}

type PostgresConfig struct {
	PostgresqlHost     string
	PostgresqlPort     string
	PostgresqlUser     string
	PostgresqlPassword string
	PostgresqlDbname   string
	PostgresqlSSLMode  bool
	PgDriver           string
}

type RabbitConfig struct {
	RabbitHost     string
	RabbitPort     string
	RabbitUser     string
	RabbitPassword string
}
