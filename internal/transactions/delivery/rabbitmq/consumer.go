package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/NikitaTsaralov/bankingApp/config"
	"github.com/NikitaTsaralov/bankingApp/internal/models"
	"github.com/NikitaTsaralov/bankingApp/internal/transactions"
	"github.com/NikitaTsaralov/bankingApp/pkg/rabbitmq"
	"github.com/NikitaTsaralov/bankingApp/pkg/rabbitmq/reconnect"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type TransactionConsumer struct {
	cfg            *config.Config
	rabbitmqChan   *reconnect.Channel
	transactionsUC transactions.UseCase
	logger         *log.Logger
}

func InitTransactionConsumer(cfg *config.Config, transactionsUC transactions.UseCase, logger *log.Logger) (*TransactionConsumer, error) {
	mqConn, err := rabbitmq.Init(cfg)
	if err != nil {
		return nil, fmt.Errorf("problem init Connection: %v", err)
	}
	rabbitmqChan, err := mqConn.Channel()
	if err != nil {
		return nil, fmt.Errorf("problem creating Channel: %v", err)
	}

	_, err = rabbitmqChan.QueueDeclare(
		cfg.Server.QueueOut, // name
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("problem rabbitmqChan.QueueDeclare: %v", err)
	}

	err = rabbitmqChan.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, errors.Wrap(err, "error  ch.Qos")
	}

	return &TransactionConsumer{
		cfg:            cfg,
		rabbitmqChan:   rabbitmqChan,
		transactionsUC: transactionsUC,
		logger:         logger,
	}, nil
}

func (c *TransactionConsumer) worker(ch *reconnect.Channel, messages <-chan amqp.Delivery) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.cfg.Server.CtxTimeoutBroker)*time.Second)
	defer cancel()

	for d := range messages {
		var transaction models.ResponseTransaction
		err := json.Unmarshal(d.Body, &transaction)
		if err != nil {
			log.Printf(" json.Unmarshal error: %v", err)
			d.Nack(false, false)
		}
		log.Printf("Unmarshal message: %s", d.Body)

		gormTransaction, err := c.transactionsUC.MoneyOperation(&transaction)
		if err != nil {
			log.Printf("creating transaction error: %v", err)
			d.Nack(false, false)
		}

		jsonBytes, err := json.Marshal(gormTransaction)
		if err != nil {
			log.Printf("json.Marshal error: %v", err)
			// d.Nack(false, false)
		}

		err = ch.PublishWithContext(ctx,
			"",        // exchange
			d.ReplyTo, // routing key
			false,     // mandatory
			false,     // immediate
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: d.CorrelationId,
				Body:          jsonBytes,
			})
		if err != nil {
			log.Printf("ch.PublishWithContext error: %v", err)
			// d.Nack(false, false)
		}

		d.Ack(false)
	}
}

func (c *TransactionConsumer) StartConsumer() error {
	msgs, err := c.rabbitmqChan.Consume(
		c.cfg.Server.QueueOut, // queue
		"",                    // consumer
		false,                 // auto-ack
		false,                 // exclusive
		false,                 // no-local
		false,                 // no-wait
		nil,                   // args
	)
	if err != nil {
		return errors.Wrap(err, "Consume")
	}

	go c.worker(c.rabbitmqChan, msgs)

	chanErr := <-c.rabbitmqChan.NotifyClose(make(chan *amqp.Error))
	log.Printf("ch.NotifyClose: %v", chanErr)
	return chanErr
}
