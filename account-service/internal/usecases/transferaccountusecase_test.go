package usecases

import (
	"errors"
	"testing"

	"github.com/matheus-oliveira-andrade/bank-statement/account-service/internal/domain"
	usecases_mock "github.com/matheus-oliveira-andrade/bank-statement/account-service/internal/usecases/mocks"
	"github.com/stretchr/testify/assert"
)

func TestTransferAccountUseCase_Handle_ErrorGettingFromAccount(t *testing.T) {
	// arrange
	mockRepo := new(usecases_mock.MockAccountRepository)
	useCase := NewTransferAccountUseCase(mockRepo)

	mockRepo.On("GetAccountByNumber", "123").Return((*domain.Account)(nil), errors.New("generic error"))

	// act
	err := useCase.Handle("123", "456", 100)

	// assert
	assert.Error(t, err)
	assert.Equal(t, "generic error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestTransferAccountUseCase_Handle_FromAccountNotFound(t *testing.T) {
	// arrange
	mockRepo := new(usecases_mock.MockAccountRepository)
	useCase := NewTransferAccountUseCase(mockRepo)

	mockRepo.On("GetAccountByNumber", "123").Return((*domain.Account)(nil), nil)

	// act
	err := useCase.Handle("123", "456", 100)

	// assert
	assert.Error(t, err)
	assert.Equal(t, "from account not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestTransferAccountUseCase_Handle_ErrorGettingToAccount(t *testing.T) {
	// arrange
	mockRepo := new(usecases_mock.MockAccountRepository)
	useCase := NewTransferAccountUseCase(mockRepo)

	fromAcc := domain.NewAccount("123", "01234567890", "John Doe")
	mockRepo.On("GetAccountByNumber", "123").Return(fromAcc, nil)
	mockRepo.On("GetAccountByNumber", "456").Return((*domain.Account)(nil), errors.New("generic error"))

	// act
	err := useCase.Handle("123", "456", 100)

	// assert
	assert.Error(t, err)
	assert.Equal(t, "generic error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestTransferAccountUseCase_Handle_ToAccountNotFound(t *testing.T) {
	// arrange
	mockRepo := new(usecases_mock.MockAccountRepository)
	useCase := NewTransferAccountUseCase(mockRepo)

	fromAcc := domain.NewAccount("123", "01234567890", "John Doe")
	mockRepo.On("GetAccountByNumber", "123").Return(fromAcc, nil)
	mockRepo.On("GetAccountByNumber", "456").Return((*domain.Account)(nil), nil)

	// act
	err := useCase.Handle("123", "456", 100)

	// assert
	assert.Error(t, err)
	assert.Equal(t, "to account not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestTransferAccountUseCase_Handle_ErrorUpdatingToAccountBalance(t *testing.T) {
	// arrange
	mockRepo := new(usecases_mock.MockAccountRepository)
	useCase := NewTransferAccountUseCase(mockRepo)

	fromAcc := domain.NewAccount("123", "01234567890", "John Doe")
	toAcc := domain.NewAccount("456", "09876543210", "Jane Doe")
	mockRepo.On("GetAccountByNumber", "123").Return(fromAcc, nil)
	mockRepo.On("GetAccountByNumber", "456").Return(toAcc, nil)

	mockRepo.On("UpdateAccountBalance", toAcc).Return(errors.New("update error"))

	// act
	err := useCase.Handle("123", "456", 100)

	// assert
	assert.Error(t, err)
	assert.Equal(t, "update error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestTransferAccountUseCase_Handle_ErrorUpdatingFromAccountBalance(t *testing.T) {
	// arrange
	mockRepo := new(usecases_mock.MockAccountRepository)
	useCase := NewTransferAccountUseCase(mockRepo)

	fromAcc := domain.NewAccount("123", "01234567890", "John Doe")
	toAcc := domain.NewAccount("456", "09876543210", "Jane Doe")
	mockRepo.On("GetAccountByNumber", "123").Return(fromAcc, nil)
	mockRepo.On("GetAccountByNumber", "456").Return(toAcc, nil)

	mockRepo.On("UpdateAccountBalance", toAcc).Return(nil)
	mockRepo.On("UpdateAccountBalance", fromAcc).Return(errors.New("update error"))

	// act
	err := useCase.Handle("123", "456", 100)

	// assert
	assert.Error(t, err)
	assert.Equal(t, "update error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestTransferAccountUseCase_Handle_Success(t *testing.T) {
	// arrange
	mockRepo := new(usecases_mock.MockAccountRepository)
	useCase := NewTransferAccountUseCase(mockRepo)

	fromAcc := domain.NewAccount("123", "01234567890", "John Doe")
	toAcc := domain.NewAccount("456", "09876543210", "Jane Doe")
	mockRepo.On("GetAccountByNumber", "123").Return(fromAcc, nil)
	mockRepo.On("GetAccountByNumber", "456").Return(toAcc, nil)

	mockRepo.On("UpdateAccountBalance", toAcc).Return(nil)
	mockRepo.On("UpdateAccountBalance", fromAcc).Return(nil)

	// act
	err := useCase.Handle("123", "456", 100)

	// assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
