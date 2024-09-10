package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/configs"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/eventhandlers"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/infrastructure/broker"
	documentgenerator "github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/infrastructure/documentgenerator"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/infrastructure/templatecompiler"
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
			case events.FundsDepositedEventKey:
				eventFundsDepositedConsume(EventPublish, dbConnection)
			case events.TransferRealizedEventKey:
				eventTransferRealizedConsume(EventPublish, dbConnection)
			case events.TransferReceivedEventKey:
				eventTransferReceivedConsume(EventPublish, dbConnection)
			case events.StatementGenerationRequestedEventKey:
				eventStatementGenerationRequested(EventPublish, dbConnection)
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
		slog.Error("error decoding event", "type", EventPublish.Type, "error", err)
		return
	}

	repository := repositories.NewAccountRepository(dbConnection)
	handler := eventhandlers.NewAccountCreatedHandler(repository)

	handler.Handler(obj)
}

func eventFundsDepositedConsume(EventPublish events.EventPublish, dbConnection *sql.DB) {
	var obj events.FundsDeposited
	err := decodeEvent([]byte(EventPublish.Data), &obj)
	if err != nil {
		slog.Error("error decoding event", "Type", EventPublish.Type, "error", err)
		return
	}

	accountRepository := repositories.NewAccountRepository(dbConnection)
	movementRepository := repositories.NewMovementRepository(dbConnection)

	handler := eventhandlers.NewFundsDepositedHandler(accountRepository, movementRepository)

	handler.Handler(obj)
}

func eventTransferRealizedConsume(EventPublish events.EventPublish, dbConnection *sql.DB) {
	var obj events.TransferRealized
	err := decodeEvent([]byte(EventPublish.Data), &obj)
	if err != nil {
		slog.Error("error decoding event", "Type", EventPublish.Type, "error", err)
		return
	}

	accountRepository := repositories.NewAccountRepository(dbConnection)
	movementRepository := repositories.NewMovementRepository(dbConnection)

	handler := eventhandlers.NewTransferRealizedHandler(accountRepository, movementRepository)

	handler.Handler(obj)
}

func eventTransferReceivedConsume(EventPublish events.EventPublish, dbConnection *sql.DB) {
	var obj events.TransferReceived
	err := decodeEvent([]byte(EventPublish.Data), &obj)
	if err != nil {
		slog.Error("error decoding event", "Type", EventPublish.Type, "error", err)
		return
	}

	accountRepository := repositories.NewAccountRepository(dbConnection)
	movementRepository := repositories.NewMovementRepository(dbConnection)

	handler := eventhandlers.NewTransferReceivedHandler(accountRepository, movementRepository)

	handler.Handler(obj)
}

func eventStatementGenerationRequested(EventPublish events.EventPublish, dbConnection *sql.DB) {
	var obj events.StatementGenerationRequested
	err := decodeEvent([]byte(EventPublish.Data), &obj)
	if err != nil {
		slog.Error("error decoding event", "Type", EventPublish.Type, "error", err)
		return
	}

	accountRepository := repositories.NewAccountRepository(dbConnection)
	statementGenerationRepository := repositories.NewStatementGenerationRepository(dbConnection)
	movementRepository := repositories.NewMovementRepository(dbConnection)
	documentGenerationApi := documentgenerator.NewGenerateDocumentApi(http.Client{})
	templateCompiler := templatecompiler.NewTemplateCompile()

	handler := eventhandlers.NewStatementGenerationRequestedHandler(
		accountRepository,
		statementGenerationRepository,
		movementRepository,
		documentGenerationApi,
		templateCompiler)

	handler.Handle(obj)
}

func decodeEvent(data []byte, obj any) error {
	return json.NewDecoder(bytes.NewReader(data)).Decode(obj)
}
