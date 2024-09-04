package usecases

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/domain"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/repositories"
)

type TriggerStatementGenerationUseCaseInterface interface {
	Handle(accountNumber string) (string, error)
}

type TriggerStatementGenerationUseCase struct {
	statementGenerationRepository repositories.StatementGenerationRepositoryInterface
	accountRepository             repositories.AccountRepositoryInterface
}

func NewTriggerStatementGenerationUseCase(
	statementGenerationRepository repositories.StatementGenerationRepositoryInterface,
	accountRepositoryInterface repositories.AccountRepositoryInterface) *TriggerStatementGenerationUseCase {
	return &TriggerStatementGenerationUseCase{
		statementGenerationRepository: statementGenerationRepository,
		accountRepository:             accountRepositoryInterface,
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

	slog.Info("statement generation created", "accountNumber", accountNumber, "triggerId", triggerId)
	return triggerId, nil
}
