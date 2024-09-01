package eventhandlers

import (
	"errors"
	"testing"

	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/domain"
	handlersmock "github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/eventhandlers/mocks"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/shared/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFundsDepositedHandler_Handler_ErrorGettingAccount(t *testing.T) {
	// arrange
	accountrepomock := new(handlersmock.MockAccountRepository)
	movementRepoMock := new(handlersmock.MockMovementRepository)
	handler := NewFundsDepositedHandler(accountrepomock, movementRepoMock)

	event := events.FundsDeposited{
		Number: "1234567890",
		Value:  100,
	}

	accountrepomock.
		On("GetAccountByNumber", event.Number).
		Return((*domain.Account)(nil), errors.New("generic error"))

	movementRepoMock.On("CreateMovement", mock.Anything).Return(nil)

	// act
	handler.Handler(event)

	// assert
	accountrepomock.AssertExpectations(t)
}

func TestFundsDepositedHandler_Handler_AccountNotFound(t *testing.T) {
	// arrange
	accountrepomock := new(handlersmock.MockAccountRepository)
	movementRepoMock := new(handlersmock.MockMovementRepository)
	handler := NewFundsDepositedHandler(accountrepomock, movementRepoMock)

	event := events.FundsDeposited{
		Number: "1234567890",
		Value:  100,
	}

	accountrepomock.
		On("GetAccountByNumber", event.Number).
		Return((*domain.Account)(nil), nil)

	movementRepoMock.On("CreateMovement", mock.Anything).Return(nil)

	// act
	handler.Handler(event)

	// assert
	accountrepomock.AssertExpectations(t)
}

func TestFundsDepositedHandler_Handler_ErrorUpdatingAccountBalance(t *testing.T) {
	// arrange
	accountrepomock := new(handlersmock.MockAccountRepository)
	movementRepoMock := new(handlersmock.MockMovementRepository)
	handler := NewFundsDepositedHandler(accountrepomock, movementRepoMock)

	acc := domain.NewAccount("1234567890", "01234567890", "John Doe")
	event := events.FundsDeposited{
		Number: "1234567890",
		Value:  100,
	}

	accountrepomock.On("GetAccountByNumber", event.Number).Return(acc, nil)
	accountrepomock.On("UpdateAccountBalance", mock.Anything).Return(errors.New("update error"))

	movementRepoMock.On("CreateMovement", mock.Anything).Return(nil)

	// act
	handler.Handler(event)

	// assert
	assert.Equal(t, int64(100), acc.Balance)
	accountrepomock.AssertExpectations(t)
}

func TestFundsDepositedHandler_Handler_Success(t *testing.T) {
	// arrange
	accountrepomock := new(handlersmock.MockAccountRepository)
	movementRepoMock := new(handlersmock.MockMovementRepository)
	handler := NewFundsDepositedHandler(accountrepomock, movementRepoMock)

	acc := domain.NewAccount("1234567890", "01234567890", "John Doe")
	acc.Balance = 100

	event := events.FundsDeposited{
		Number: "1234567890",
		Value:  100,
	}

	accountrepomock.On("GetAccountByNumber", event.Number).Return(acc, nil)
	accountrepomock.On("UpdateAccountBalance", acc).Return(nil)

	movementRepoMock.On("CreateMovement", mock.Anything).Return(nil)

	// act
	handler.Handler(event)

	// assert
	assert.Equal(t, int64(200), acc.Balance)
	accountrepomock.AssertExpectations(t)
}
