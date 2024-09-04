package usecases_mocks

import (
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) GetAccountByNumber(number string) (*domain.Account, error) {
	args := m.Called(number)
	return args.Get(0).(*domain.Account), args.Error(1)
}

func (m *MockAccountRepository) CreateAccount(account *domain.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *MockAccountRepository) UpdateAccountBalance(account *domain.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

type MockMovementRepository struct {
	mock.Mock
}

func (m *MockMovementRepository) CreateMovement(movement *domain.Movement) error {
	args := m.Called(movement)
	return args.Error(0)
}
