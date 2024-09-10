package handlersmock

import (
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockStatementGenerationRepository struct {
	mock.Mock
}

func (m *MockStatementGenerationRepository) CreateStatementGeneration(statementGeneration *domain.StatementGeneration) (string, error) {
	args := m.Called(statementGeneration)
	return args.String(0), args.Error(1)
}

func (m *MockStatementGenerationRepository) HasStatementGenerationRunning(accountNumber string) (bool, error) {
	args := m.Called(accountNumber)
	return args.Bool(0), args.Error(1)
}

func (m *MockStatementGenerationRepository) GetStatementGeneration(accountNumber string) (*domain.StatementGeneration, error) {
	args := m.Called(accountNumber)
	return args.Get(0).(*domain.StatementGeneration), args.Error(1)
}

func (m *MockStatementGenerationRepository) UpdateStatementGeneration(statementGeneration *domain.StatementGeneration) error {
	args := m.Called(statementGeneration)
	return args.Error(0)
}
