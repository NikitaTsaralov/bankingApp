package rabbitmq

import (
	"fmt"

	"github.com/NikitaTsaralov/bankingApp/config"
	"github.com/NikitaTsaralov/bankingApp/pkg/rabbitmq/reconnect"
)

func Init(cfg *config.Config) (*reconnect.Connection, error) {
	connAddr := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.Rabbit.RabbitUser,
		cfg.Rabbit.RabbitPassword,
		cfg.Rabbit.RabbitHost,
		cfg.Rabbit.RabbitPort,
	)

	return reconnect.Dial(connAddr)
}
