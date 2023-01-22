package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/NikitaTsaralov/bankingApp/config"
	"github.com/NikitaTsaralov/bankingApp/internal/models"
	"github.com/NikitaTsaralov/bankingApp/pkg/rabbitmq"
	"github.com/NikitaTsaralov/bankingApp/pkg/rabbitmq/reconnect"
	"github.com/NikitaTsaralov/bankingApp/pkg/utils"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type TransactionPublisher struct {
	rabbitmqChan *reconnect.Channel
	cfg          *config.Config
	logger       *log.Logger
}

func InitTransactionPublisher(cfg *config.Config, logger *log.Logger) (*TransactionPublisher, error) {
	mqConn, err := rabbitmq.Init(cfg)
	if err != nil {
		return nil, fmt.Errorf("problem init Connection: %v", err)
	}
	rabbitmqChan, err := mqConn.Channel()
	if err != nil {
		return nil, fmt.Errorf("problem creating Channel: %v", err)
	}

	return &TransactionPublisher{
		cfg:          cfg,
		rabbitmqChan: rabbitmqChan,
		logger:       logger,
	}, nil
}

func (publisher *TransactionPublisher) Setup() error {
	// queue
	return nil
}

func (publisher *TransactionPublisher) Close() {
	if err := publisher.rabbitmqChan.Close(); err != nil {
		log.Printf("error TransactionPublisher CloseChan: %v", err)
	}
}

func (publisher *TransactionPublisher) Publish(transaction *models.ResponseTransaction) (res *models.ResponseTransaction, err error) {
	queue, err := publisher.rabbitmqChan.QueueDeclare(
		publisher.cfg.Server.QueueIn,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "error ch.QueueDeclare")
	}

	msgs, err := publisher.rabbitmqChan.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return nil, errors.Wrap(err, "error ch.Consume")
	}

	corrId := utils.RandomString(32)

	jsonBytes, err := json.Marshal(transaction)
	if err != nil {
		return nil, errors.Wrap(err, "error json.Unmarshal")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(publisher.cfg.Server.CtxTimeoutBroker)*time.Second)
	defer cancel()

	err = publisher.rabbitmqChan.PublishWithContext(ctx,
		"",                            // exchange
		publisher.cfg.Server.QueueOut, // routing key
		false,                         // mandatory
		false,                         // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       queue.Name,
			Body:          jsonBytes,
		})
	if err != nil {
		return nil, errors.Wrap(err, "error json.Unmarshal")
	}

	for d := range msgs {
		if corrId == d.CorrelationId {
			err := json.Unmarshal(d.Body, &res)
			if err != nil {
				return nil, errors.Wrap(err, "error json.Unmarshal")
			}
			d.Ack(false)
			return res, nil
		}
	}

	return nil, nil
}
