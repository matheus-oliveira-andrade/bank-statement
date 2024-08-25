package usecases_mock

import (
	"github.com/matheus-oliveira-andrade/bank-statement/account-service/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) GetAccountByNumber(number string) (*domain.Account, error) {
	args := m.Called(number)
	return args.Get(0).(*domain.Account), args.Error(1)
}

func (m *MockAccountRepository) GetAccountByDocument(document string) (*domain.Account, error) {
	args := m.Called(document)
	return args.Get(0).(*domain.Account), args.Error(1)
}

func (m *MockAccountRepository) GetNextAccountNumber() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockAccountRepository) CreateAccount(account *domain.Account) (string, error) {
	args := m.Called(account)
	return args.String(0), args.Error(1)
}

func (m *MockAccountRepository) UpdateAccountBalance(account *domain.Account) error {
	args := m.Called(account)
	return args.Error(0)
}
