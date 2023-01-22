package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/NikitaTsaralov/bankingApp/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient struct {
	cfg    *config.Config
	logger *log.Logger
}

func Init(cfg *config.Config, logger *log.Logger) *RabbitMQClient {
	return &RabbitMQClient{
		cfg: cfg,
	}
}

func (rabbit *RabbitMQClient) Send(msg []byte) (err error) {
	conn, err := amqp.Dial(
		fmt.Sprintf("amqp://%s:%s@%s:%s/",
			rabbit.cfg.Rabbit.RabbitUser,
			rabbit.cfg.Rabbit.RabbitPassword,
			rabbit.cfg.Rabbit.RabbitHost,
			rabbit.cfg.Rabbit.RabbitPort,
		))
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		rabbit.cfg.Rabbit.Queue, // name
		true,                    // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         msg,
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %v", err)
	}
	log.Printf(" [x] Sent %s", msg)
	return
}
