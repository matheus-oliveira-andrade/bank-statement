package usecases

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/domain"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/repositories"
)

type GetStatementGenerationUseCaseInterface interface {
	Handle(id string) (string, error)
}

type GetStatementGenerationUseCase struct {
	statementGenerationRepository repositories.StatementGenerationRepositoryInterface
}

func NewGetStatementGenerationUseCase(
	statementGenerationRepository repositories.StatementGenerationRepositoryInterface,
) *GetStatementGenerationUseCase {
	return &GetStatementGenerationUseCase{
		statementGenerationRepository: statementGenerationRepository,
	}
}

func (us *GetStatementGenerationUseCase) Handle(id string) (string, error) {
	sg, err := us.statementGenerationRepository.GetStatementGenerationById(id)

	if err != nil {
		slog.Info("error getting statement generation", "id", id)
		return "", fmt.Errorf("error getting statement generation")
	}

	if sg == nil {
		slog.Error("statement generation not found", "id", id)
		return "", fmt.Errorf("statement generation not found")
	}

	if sg.Status == domain.StatementGenerationRunnning {
		return "", nil
	}

	if sg.Error != "" {
		return "", errors.New(sg.Error)
	}

	return sg.DocumentContent, nil
}
