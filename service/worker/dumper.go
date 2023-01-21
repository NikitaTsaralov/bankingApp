package worker

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/NikitaTsaralov/bankingApp/config"
	"github.com/NikitaTsaralov/bankingApp/internal/models"
	"github.com/NikitaTsaralov/bankingApp/service/rabbitmq"
	"github.com/jinzhu/gorm"
)

type DBDumper struct {
	cfg      *config.Config
	database *gorm.DB
	logger   *log.Logger
}

func Init(cfg *config.Config, database *gorm.DB, logger *log.Logger) *DBDumper {
	return &DBDumper{
		cfg:      cfg,
		database: database,
		logger:   logger,
	}
}

func (dumper *DBDumper) moneyOperation(transaction *models.ResponseTransaction) (err error) {
	if dbc := dumper.database.Create(&models.Transaction{
		AccountId: transaction.AccountId,
		Amount:    transaction.Amount,
	}); dbc.Error != nil {
		return fmt.Errorf("Error creating Transaction: %v", dbc.Error)
	}

	gormAccount := &models.Account{}
	if dbc := dumper.database.Where("id = ? ", transaction.AccountId).First(&gormAccount); dbc.Error != nil {
		return fmt.Errorf("Error getting Account: %v", dbc.Error)
	}

	gormAccount.Balance += transaction.Amount
	if dbc := dumper.database.Save(&gormAccount); dbc.Error != nil {
		return fmt.Errorf("Error updating Account: %v", dbc.Error)
	}
	return
}

func (dumper *DBDumper) Run() error {
	conn, err := rabbitmq.Dial(
		fmt.Sprintf("amqp://%s:%s@%s:%s/",
			dumper.cfg.Rabbit.RabbitUser,
			dumper.cfg.Rabbit.RabbitPassword,
			dumper.cfg.Rabbit.RabbitHost,
			dumper.cfg.Rabbit.RabbitPort,
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
		dumper.cfg.Rabbit.Queue, // name
		true,                    // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %v", err)
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return fmt.Errorf("Failed to set QoS: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("Failed to register a consumer: %v", err)
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			var transaction models.ResponseTransaction

			err := json.Unmarshal(d.Body, &transaction)
			if err != nil {
				log.Printf("parse transaction error: %v", err)
				d.Nack(false, false)
			}
			log.Printf("Unmarshal message: %s", d.Body)

			err = dumper.moneyOperation(&transaction)
			if err != nil {
				log.Printf("creating transaction error: %v", err)
				d.Nack(false, false)
			}

			log.Printf("Done")
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	return nil
}
