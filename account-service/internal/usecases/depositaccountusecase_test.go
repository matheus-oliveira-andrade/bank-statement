package usecases

import (
	"errors"
	"testing"

	"github.com/matheus-oliveira-andrade/bank-statement/account-service/internal/domain"
	usecases_mock "github.com/matheus-oliveira-andrade/bank-statement/account-service/internal/usecases/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDepositAccountUseCase_Handle_ErrorGettingAccount(t *testing.T) {
	// act
	mockRepo := new(usecases_mock.MockAccountRepository)
	useCase := NewDepositAccountUseCase(mockRepo)

	acc := domain.NewAccount("4", "01234567890", "John Doo")

	mockRepo.On("GetAccountByNumber", acc.Number).Return(acc, errors.New("generic error"))

	// act
	err := useCase.Handle(acc.Number, 150)

	// assert
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestDepositAccountUseCase_Handle_AccountNotFound(t *testing.T) {
	// act
	mockRepo := new(usecases_mock.MockAccountRepository)
	useCase := NewDepositAccountUseCase(mockRepo)

	acc := domain.NewAccount("4", "01234567890", "John Doo")

	mockRepo.On("GetAccountByNumber", acc.Number).Return((*domain.Account)(nil), nil)

	// act
	err := useCase.Handle(acc.Number, 150)

	// assert
	assert.Error(t, err)
	assert.Equal(t, "account not found", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestDepositAccountUseCase_Handle_Success(t *testing.T) {
	// act
	mockRepo := new(usecases_mock.MockAccountRepository)
	useCase := NewDepositAccountUseCase(mockRepo)

	acc := domain.NewAccount("4", "01234567890", "John Doo")

	mockRepo.On("GetAccountByNumber", acc.Number).Return(acc, nil)
	mockRepo.On("UpdateAccountBalance", mock.Anything).Return(nil)

	// act
	err := useCase.Handle(acc.Number, 150)

	// assert
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}
