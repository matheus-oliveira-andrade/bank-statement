package usecases

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/matheus-oliveira-andrade/bank-statement/account-service/internal/domain"
	usecases_mock "github.com/matheus-oliveira-andrade/bank-statement/account-service/internal/usecases/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDepositAccountUseCase_Handle_ErrorGettingAccount(t *testing.T) {
	// act
	mockRepo := new(usecases_mock.MockAccountRepository)
	mockBroker := new(usecases_mock.MockBroker)
	mockIdempotencyRepository := new(usecases_mock.MockIdempotencyRepository)

	useCase := NewDepositAccountUseCase(mockRepo, mockBroker, mockIdempotencyRepository)

	acc := domain.NewAccount("4", "01234567890", "John Doo")

	mockRepo.On("GetAccountByNumber", acc.Number).Return(acc, errors.New("generic error"))
	mockBroker.On("Produce", mock.Anything, mock.Anything).Return(nil)

	idempotencyKey, _ := uuid.NewUUID()

	mockIdempotencyRepository.On("HasKey", idempotencyKey.String()).Return(false, nil)
	mockIdempotencyRepository.On("CreateKey", idempotencyKey.String()).Return(nil)

	// act
	err := useCase.Handle(acc.Number, 150, idempotencyKey.String())

	// assert
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestDepositAccountUseCase_Handle_AccountNotFound(t *testing.T) {
	// act
	mockRepo := new(usecases_mock.MockAccountRepository)
	mockBroker := new(usecases_mock.MockBroker)
	mockIdempotencyRepository := new(usecases_mock.MockIdempotencyRepository)

	useCase := NewDepositAccountUseCase(mockRepo, mockBroker, mockIdempotencyRepository)

	acc := domain.NewAccount("4", "01234567890", "John Doo")

	mockRepo.On("GetAccountByNumber", acc.Number).Return((*domain.Account)(nil), nil)
	mockBroker.On("Produce", mock.Anything, mock.Anything).Return(nil)

	idempotencyKey, _ := uuid.NewUUID()

	mockIdempotencyRepository.On("HasKey", idempotencyKey.String()).Return(false, nil)
	mockIdempotencyRepository.On("CreateKey", idempotencyKey.String()).Return(nil)

	// act
	err := useCase.Handle(acc.Number, 150, idempotencyKey.String())

	// assert
	assert.Error(t, err)
	assert.Equal(t, "account not found", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestDepositAccountUseCase_Handle_Success(t *testing.T) {
	// act
	mockRepo := new(usecases_mock.MockAccountRepository)
	mockBroker := new(usecases_mock.MockBroker)
	mockIdempotencyRepository := new(usecases_mock.MockIdempotencyRepository)

	useCase := NewDepositAccountUseCase(mockRepo, mockBroker, mockIdempotencyRepository)

	acc := domain.NewAccount("4", "01234567890", "John Doo")

	mockRepo.On("GetAccountByNumber", acc.Number).Return(acc, nil)
	mockRepo.On("UpdateAccountBalance", mock.Anything).Return(nil)
	mockBroker.On("Produce", mock.Anything, mock.Anything).Return(nil)

	idempotencyKey, _ := uuid.NewUUID()

	mockIdempotencyRepository.On("HasKey", idempotencyKey.String()).Return(false, nil)
	mockIdempotencyRepository.On("CreateKey", idempotencyKey.String()).Return(nil)

	// act
	err := useCase.Handle(acc.Number, 150, idempotencyKey.String())

	// assert
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockIdempotencyRepository.AssertExpectations(t)
}

func TestDepositAccountUseCase_Handle_IdempotencyKeyAlreadyUsed(t *testing.T) {
	// act
	mockRepo := new(usecases_mock.MockAccountRepository)
	mockBroker := new(usecases_mock.MockBroker)
	mockIdempotencyRepository := new(usecases_mock.MockIdempotencyRepository)

	useCase := NewDepositAccountUseCase(mockRepo, mockBroker, mockIdempotencyRepository)

	acc := domain.NewAccount("4", "01234567890", "John Doo")

	idempotencyKey, _ := uuid.NewUUID()

	mockIdempotencyRepository.On("HasKey", idempotencyKey.String()).Return(true, nil)

	// act
	err := useCase.Handle(acc.Number, 150, idempotencyKey.String())

	// assert
	assert.Error(t, err)
	assert.Equal(t, err, errors.New("idempotency key already processed"))

	mockRepo.AssertExpectations(t)
	mockBroker.AssertExpectations(t)
	mockIdempotencyRepository.AssertExpectations(t)
}

func TestDepositAccountUseCase_Handle_HasKeyErrorRetriving(t *testing.T) {
	// act
	mockRepo := new(usecases_mock.MockAccountRepository)
	mockBroker := new(usecases_mock.MockBroker)
	mockIdempotencyRepository := new(usecases_mock.MockIdempotencyRepository)

	useCase := NewDepositAccountUseCase(mockRepo, mockBroker, mockIdempotencyRepository)

	acc := domain.NewAccount("4", "01234567890", "John Doo")

	idempotencyKey, _ := uuid.NewUUID()

	genericError := errors.New("generic error")
	mockIdempotencyRepository.On("HasKey", idempotencyKey.String()).Return(false, genericError)

	// act
	err := useCase.Handle(acc.Number, 150, idempotencyKey.String())

	// assert
	assert.Error(t, err)
	assert.Equal(t, err, genericError)

	mockRepo.AssertExpectations(t)
	mockBroker.AssertExpectations(t)
	mockIdempotencyRepository.AssertExpectations(t)
}

func TestDepositAccountUseCase_Handle_SuccessButErrorSavingIdempotencyKeyUsed(t *testing.T) {
	// act
	mockRepo := new(usecases_mock.MockAccountRepository)
	mockBroker := new(usecases_mock.MockBroker)
	mockIdempotencyRepository := new(usecases_mock.MockIdempotencyRepository)

	useCase := NewDepositAccountUseCase(mockRepo, mockBroker, mockIdempotencyRepository)

	acc := domain.NewAccount("4", "01234567890", "John Doo")

	mockRepo.On("GetAccountByNumber", acc.Number).Return(acc, nil)
	mockRepo.On("UpdateAccountBalance", mock.Anything).Return(nil)
	mockBroker.On("Produce", mock.Anything, mock.Anything).Return(nil)

	idempotencyKey, _ := uuid.NewUUID()

	mockIdempotencyRepository.On("HasKey", idempotencyKey.String()).Return(false, nil)

	errorSavingUsedIdempotencyKey := errors.New("error saving used idempotency key")
	mockIdempotencyRepository.On("CreateKey", idempotencyKey.String()).Return(errorSavingUsedIdempotencyKey)

	// act
	err := useCase.Handle(acc.Number, 150, idempotencyKey.String())

	// assert
	assert.Error(t, err)
	assert.Equal(t, err, errorSavingUsedIdempotencyKey)

	mockRepo.AssertExpectations(t)
	mockIdempotencyRepository.AssertExpectations(t)
}
