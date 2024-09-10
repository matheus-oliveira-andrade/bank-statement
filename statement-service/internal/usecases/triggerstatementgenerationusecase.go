package usecases

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/domain"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/infrastructure/broker"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/repositories"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/shared/events"
)

type TriggerStatementGenerationUseCaseInterface interface {
	Handle(accountNumber string) (string, error)
}

type TriggerStatementGenerationUseCase struct {
	statementGenerationRepository repositories.StatementGenerationRepositoryInterface
	accountRepository             repositories.AccountRepositoryInterface
	broker                        broker.BrokerInterface
}

func NewTriggerStatementGenerationUseCase(
	statementGenerationRepository repositories.StatementGenerationRepositoryInterface,
	accountRepositoryInterface repositories.AccountRepositoryInterface,
	broker broker.BrokerInterface,
) *TriggerStatementGenerationUseCase {
	return &TriggerStatementGenerationUseCase{
		statementGenerationRepository: statementGenerationRepository,
		accountRepository:             accountRepositoryInterface,
		broker:                        broker,
	}
}

func (us *TriggerStatementGenerationUseCase) Handle(accountNumber string) (string, error) {
	if acc, _ := us.accountRepository.GetAccountByNumber(accountNumber); acc == nil {
		slog.Info("account not found", "accountNumber", accountNumber)
		return "", fmt.Errorf("account not found: %v", accountNumber)
	}

	hasGenerationRunning, err := us.statementGenerationRepository.HasStatementGenerationRunning(accountNumber)
	if err != nil {
		slog.Info(err.Error(), "accountNumber", accountNumber)
		return "", errors.New("error search for statement generation running")
	}

	if hasGenerationRunning {
		slog.Info("already have a statement generation running", "accountNumber", accountNumber)
		return "", errors.New("already have a statement generation running")
	}

	statementGeneration, err := domain.NewStatementGeneration(accountNumber)
	if err != nil {
		slog.Info("Error creating statement generation", "err", err)
		return "", err
	}

	triggerId, err := us.statementGenerationRepository.CreateStatementGeneration(statementGeneration)
	if err != nil {
		slog.Info(err.Error(), "accountNumber", accountNumber)
		return "", errors.New("error creating statement generation")
	}

	event := events.NewStatementGenerationRequested(triggerId, accountNumber)
	eventPublish, err := events.NewEventPublish(event)
	if err != nil {
		slog.Error("error creating event publish", "event", event)
		return "", err
	}

	us.broker.Produce(eventPublish, &broker.ProduceConfigs{Topic: "statement"})

	slog.Info("statement generation created", "accountNumber", accountNumber, "triggerId", triggerId)
	return triggerId, nil
}
