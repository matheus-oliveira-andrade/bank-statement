package eventhandlers

import (
	"testing"

	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/domain"
	handlersmock "github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/eventhandlers/mocks"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/shared/events"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewAccountCreatedHandler(t *testing.T) {
	// Arrange
	accountRepoMock := new(handlersmock.MockAccountRepository)

	// Act
	handler := NewAccountCreatedHandler(accountRepoMock)

	// Assert
	require.NotNil(t, handler)
}

func TestAccountCreatedHandler_Handler(t *testing.T) {
	// Arrange
	accountRepoMock := new(handlersmock.MockAccountRepository)
	handler := NewAccountCreatedHandler(accountRepoMock)

	event := events.AccountCreated{
		Number:   "123456789",
		Document: "12345678900",
		Name:     "John Doe",
	}

	expectedAccount := domain.NewAccount(event.Number, event.Document, event.Name)

	accountRepoMock.
		On("CreateAccount", mock.AnythingOfType("*domain.Account")).
		Return(nil).
		Run(func(args mock.Arguments) {
			acc := args.Get(0).(*domain.Account)
			require.Equal(t, expectedAccount.Number, acc.Number)
			require.Equal(t, expectedAccount.Document, acc.Document)
			require.Equal(t, expectedAccount.Name, acc.Name)
		})

	// Act
	handler.Handler(event)

	// Assert
	accountRepoMock.AssertExpectations(t)
}
