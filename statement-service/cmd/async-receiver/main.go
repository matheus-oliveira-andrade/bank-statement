package main

import (
	"bytes"
	"database/sql"
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

	consumedMessages, err := ch.Consume(
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
		for consumedMessage := range consumedMessages {
			var EventPublish events.EventPublish
			err := decodeEvent(consumedMessage.Body, &EventPublish)
			if err != nil {
				fmt.Println(err)
			}

			switch EventPublish.Type {
			case events.AccountCreatedEventKey:
				eventAccountCreatedConsume(EventPublish, dbConnection)
			default:
				slog.Info("event type not mapped", "eventType", EventPublish.Type)
			}
		}
	}()

	log.Printf("[*] Waiting for events. To exit press CTRL+C")
	<-forever
}

func eventAccountCreatedConsume(EventPublish events.EventPublish, dbConnection *sql.DB) {
	var obj events.AccountCreated
	err := decodeEvent([]byte(EventPublish.Data), &obj)
	if err != nil {
		slog.Error("error decoding event", "error", err)
		return
	}

	repository := repositories.NewAccountRepository(dbConnection)
	handler := eventhandlers.NewAccountCreatedHandler(repository)

	handler.Handler(obj)
}

func decodeEvent(data []byte, obj any) error {
	return json.NewDecoder(bytes.NewReader(data)).Decode(obj)
}
