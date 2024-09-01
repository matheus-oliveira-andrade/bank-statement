package eventhandlers

import (
	"errors"
	"testing"

	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/domain"
	handlersmock "github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/eventhandlers/mocks"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/shared/events"
	"github.com/stretchr/testify/mock"
)

func getTestEvent() events.TransferRealized {
	return events.TransferRealized{
		FromNumber: "123456789",
		ToNumber:   "987654321",
		Value:      100.0,
		Balance:    900.0,
	}
}

func getTestAccount() *domain.Account {
	return &domain.Account{
		Number:   "123456789",
		Name:     "John Doe",
		Document: "12345678901",
		Balance:  1000.0,
	}
}

func TestTransferRealizedHandler_Success(t *testing.T) {
	// Arrange
	accountRepo := new(handlersmock.MockAccountRepository)
	movementRepo := new(handlersmock.MockMovementRepository)
	handler := NewTransferRealizedHandler(accountRepo, movementRepo)

	event := getTestEvent()
	account := getTestAccount()

	accountRepo.On("GetAccountByNumber", event.FromNumber).Return(account, nil)
	accountRepo.On("UpdateAccountBalance", mock.Anything).Return(nil)
	movementRepo.On("CreateMovement", mock.Anything).Return(nil)

	// Act
	handler.Handler(event)

	// Assert
	accountRepo.AssertExpectations(t)
	movementRepo.AssertExpectations(t)
}

func TestTransferRealizedHandler_AccountNotFound(t *testing.T) {
	// Arrange
	accountRepo := new(handlersmock.MockAccountRepository)
	movementRepo := new(handlersmock.MockMovementRepository)
	handler := NewTransferRealizedHandler(accountRepo, movementRepo)

	event := getTestEvent()

	accountRepo.On("GetAccountByNumber", event.FromNumber).Return((*domain.Account)(nil), nil)

	// Act
	handler.Handler(event)

	// Assert
	accountRepo.AssertExpectations(t)
	movementRepo.AssertNotCalled(t, "CreateMovement", mock.Anything)
	accountRepo.AssertNotCalled(t, "UpdateAccountBalance", mock.Anything)
}

func TestTransferRealizedHandler_DBErrorOnGetAccount(t *testing.T) {
	// Arrange
	accountRepo := new(handlersmock.MockAccountRepository)
	movementRepo := new(handlersmock.MockMovementRepository)
	handler := NewTransferRealizedHandler(accountRepo, movementRepo)

	event := getTestEvent()

	accountRepo.On("GetAccountByNumber", event.FromNumber).Return((*domain.Account)(nil), errors.New("db error"))

	// Act
	handler.Handler(event)

	// Assert
	accountRepo.AssertExpectations(t)
	movementRepo.AssertNotCalled(t, "CreateMovement", mock.Anything)
	accountRepo.AssertNotCalled(t, "UpdateAccountBalance", mock.Anything)
}

func TestTransferRealizedHandler_DBErrorOnUpdateBalance(t *testing.T) {
	// Arrange
	accountRepo := new(handlersmock.MockAccountRepository)
	movementRepo := new(handlersmock.MockMovementRepository)
	handler := NewTransferRealizedHandler(accountRepo, movementRepo)

	event := getTestEvent()
	account := getTestAccount()

	accountRepo.On("GetAccountByNumber", event.FromNumber).Return(account, nil)
	accountRepo.On("UpdateAccountBalance", mock.Anything).Return(errors.New("db error"))

	// Act
	handler.Handler(event)

	// Assert
	accountRepo.AssertExpectations(t)
	movementRepo.AssertNotCalled(t, "CreateMovement", mock.Anything)
}

func TestTransferRealizedHandler_DBErrorOnCreateMovement(t *testing.T) {
	// Arrange
	accountRepo := new(handlersmock.MockAccountRepository)
	movementRepo := new(handlersmock.MockMovementRepository)
	handler := NewTransferRealizedHandler(accountRepo, movementRepo)

	event := getTestEvent()
	account := getTestAccount()

	accountRepo.On("GetAccountByNumber", event.FromNumber).Return(account, nil)
	accountRepo.On("UpdateAccountBalance", mock.Anything).Return(nil)
	movementRepo.On("CreateMovement", mock.Anything).Return(errors.New("db error"))

	// Act
	handler.Handler(event)

	// Assert
	accountRepo.AssertExpectations(t)
	movementRepo.AssertExpectations(t)
}
