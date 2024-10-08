package usecases

import (
	"errors"
	"testing"

	"github.com/matheus-oliveira-andrade/bank-statement/account-service/internal/domain"
	usecases_mock "github.com/matheus-oliveira-andrade/bank-statement/account-service/internal/usecases/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateAccountUseCase_Handle_Success(t *testing.T) {
	// act
	mockRepo := new(usecases_mock.MockAccountRepository)
	mockBroker := new(usecases_mock.MockBroker)
	useCase := NewCreateAccountUseCase(mockRepo, mockBroker)

	document := "12345678901"

	mockRepo.On("GetAccountByDocument", document).Return((*domain.Account)(nil), nil)
	mockRepo.On("GetNextAccountNumber").Return("987654321", nil)
	mockRepo.On("CreateAccount", mock.Anything).Return("1", nil)

	mockBroker.On("Produce", mock.Anything, mock.Anything).Return(nil)

	// act
	id, err := useCase.Handle(document, "John Doe")

	// assert
	assert.NoError(t, err)
	assert.Equal(t, "1", id)
	mockRepo.AssertExpectations(t)
}

func TestCreateAccountUseCase_Handle_DocumentInUse(t *testing.T) {
	// arange
	mockRepo := new(usecases_mock.MockAccountRepository)
	mockBroker := new(usecases_mock.MockBroker)
	useCase := NewCreateAccountUseCase(mockRepo, mockBroker)

	existingAccount := &domain.Account{
		Id:       "1",
		Number:   "987654321",
		Name:     "Jane Doe",
		Document: "12345678901",
		Balance:  1000.0,
	}
	mockRepo.On("GetAccountByDocument", "12345678901").Return(existingAccount, nil)

	mockBroker.On("Produce", mock.Anything, mock.Anything).Return(nil)
	// act
	id, err := useCase.Handle("12345678901", "John Doe")

	// assert
	assert.Error(t, err)
	assert.Equal(t, "", id)
	assert.Equal(t, "document in use by another account", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestCreateAccountUseCase_Handle_GetNextAccountNumberError(t *testing.T) {
	// arrange
	mockRepo := new(usecases_mock.MockAccountRepository)
	mockBroker := new(usecases_mock.MockBroker)
	useCase := NewCreateAccountUseCase(mockRepo, mockBroker)

	mockRepo.On("GetAccountByDocument", "12345678901").Return((*domain.Account)(nil), nil)
	mockRepo.On("GetNextAccountNumber").Return("", errors.New("error getting next account number"))

	mockBroker.On("Produce", mock.Anything, mock.Anything).Return(nil)

	// act
	id, err := useCase.Handle("12345678901", "John Doe")

	// assert
	assert.Error(t, err)
	assert.Equal(t, "", id)
	assert.Equal(t, "error getting next account number", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestCreateAccountUseCase_Handle_CreateAccountError(t *testing.T) {
	// arrange
	mockRepo := new(usecases_mock.MockAccountRepository)
	mockBroker := new(usecases_mock.MockBroker)
	useCase := NewCreateAccountUseCase(mockRepo, mockBroker)

	mockRepo.On("GetAccountByDocument", "12345678901").Return((*domain.Account)(nil), nil)
	mockRepo.On("GetNextAccountNumber").Return("987654321", nil)
	mockRepo.On("CreateAccount", mock.Anything).Return("", errors.New("error creating account"))

	mockBroker.On("Produce", mock.Anything, mock.Anything).Return(nil)

	// act
	id, err := useCase.Handle("12345678901", "John Doe")

	// assert
	assert.Error(t, err)
	assert.Equal(t, "", id)
	assert.Equal(t, "error creating account", err.Error())
	mockRepo.AssertExpectations(t)
}
