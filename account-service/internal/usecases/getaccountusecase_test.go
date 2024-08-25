package usecases

import (
	"errors"
	"testing"

	"github.com/matheus-oliveira-andrade/bank-statement/account-service/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestGetAccountUseCase_Handle_Success(t *testing.T) {
	// arrange
	acc := domain.NewAccount(
		"1",
		"01234567890",
		"John Dee",
	)

	mockAccountRepo := new(MockAccountRepository)
	mockAccountRepo.On("GetAccountByNumber", acc.Number).Return(acc, nil)

	useCase := NewGetAccountUseCase(mockAccountRepo)

	// act
	result, err := useCase.Handle(acc.Number)

	// assert
	assert.NotNil(t, result)
	assert.Nil(t, err)

	assert.Equal(t, result.Id, acc.Id)
	assert.Equal(t, result.Number, acc.Number)
	assert.Equal(t, result.Name, acc.Name)
}

func TestGetAccountUseCase_Handle_NotFound(t *testing.T) {
	// arrange
	acc := domain.NewAccount(
		"1",
		"01234567890",
		"John Dee",
	)

	mockAccountRepo := new(MockAccountRepository)
	mockAccountRepo.On("GetAccountByNumber", acc.Number).Return((*domain.Account)(nil), nil)

	useCase := NewGetAccountUseCase(mockAccountRepo)

	// act
	result, err := useCase.Handle(acc.Number)

	// assert
	assert.Nil(t, result)
	assert.Nil(t, err)
}

func TestGetAccountUseCase_Handle_Error(t *testing.T) {
	// arrange
	acc := domain.NewAccount(
		"1",
		"01234567890",
		"John Dee",
	)

	expectedError := errors.New("error getting account")

	mockAccountRepo := new(MockAccountRepository)
	mockAccountRepo.On("GetAccountByNumber", acc.Number).Return((*domain.Account)(nil), expectedError)

	useCase := NewGetAccountUseCase(mockAccountRepo)

	// act
	result, err := useCase.Handle(acc.Number)

	// assert
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
}
