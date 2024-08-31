package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"

	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/configs"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/eventhandlers"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/infrastructure/broker"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/logger"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/repositories"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/shared/events"
)

func main() {
	configs.InitConfigFile()
	logger.SetupLogger()

	brokerConn, err := broker.NewConnection(broker.BuildConnectionUrl())
	if err != nil {
		panic(err)
	}

	defer brokerConn.Close()

	dbConnection := repositories.NewDBConnection()
	if dbConnection == nil {
		panic("error connecting to database")
	}

	ch, err := brokerConn.Channel()
	if err != nil {
		panic(err)
	}

	defer ch.Close()

	msgs, err := ch.Consume(
		"statement-service-queue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			var EventPublish events.EventPublish
			err := json.NewDecoder(bytes.NewReader(d.Body)).Decode(&EventPublish)
			if err != nil {
				fmt.Println(err)
			}

			switch EventPublish.Type {
			case "AccountCreated":
				var obj events.AccountCreated

				err := json.NewDecoder(bytes.NewReader([]byte(EventPublish.Data))).Decode(&obj)
				if err != nil {
					slog.Error("error decoding event", "error", err)
				}

				repository := repositories.NewAccountRepository(dbConnection)
				handler := eventhandlers.NewAccountCreatedHandler(repository)

				handler.Handler(obj)

			default:
				slog.Info("event type not mapped", "eventType", EventPublish.Type)
			}
		}
	}()

	log.Printf(" [*] Waiting for events. To exit press CTRL+C")
	<-forever
}
