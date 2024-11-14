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

func TestTransferAccountUseCase_Handle_ErrorGettingFromAccount(t *testing.T) {
	// arrange
	mockRepo := new(usecases_mock.MockAccountRepository)
	mockBroker := new(usecases_mock.MockBroker)
	mockIdempotencyRepository := new(usecases_mock.MockIdempotencyRepository)

	useCase := NewTransferAccountUseCase(mockRepo, mockBroker, mockIdempotencyRepository)

	mockRepo.On("GetAccountByNumber", "123").Return((*domain.Account)(nil), errors.New("generic error"))

	mockBroker.On("Produce", mock.Anything, mock.Anything).Return(nil)

	idempotencyKey, _ := uuid.NewUUID()

	mockIdempotencyRepository.On("HasKey", idempotencyKey.String()).Return(false, nil)
	mockIdempotencyRepository.On("CreateKey", idempotencyKey.String()).Return(nil)

	// act
	err := useCase.Handle("123", "456", 100, idempotencyKey.String())

	// assert
	assert.Error(t, err)
	assert.Equal(t, "generic error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestTransferAccountUseCase_Handle_FromAccountNotFound(t *testing.T) {
	// arrange
	mockRepo := new(usecases_mock.MockAccountRepository)
	mockBroker := new(usecases_mock.MockBroker)
	mockIdempotencyRepository := new(usecases_mock.MockIdempotencyRepository)

	useCase := NewTransferAccountUseCase(mockRepo, mockBroker, mockIdempotencyRepository)

	mockRepo.On("GetAccountByNumber", "123").Return((*domain.Account)(nil), nil)

	mockBroker.On("Produce", mock.Anything, mock.Anything).Return(nil)

	idempotencyKey, _ := uuid.NewUUID()

	mockIdempotencyRepository.On("HasKey", idempotencyKey.String()).Return(false, nil)
	mockIdempotencyRepository.On("CreateKey", idempotencyKey.String()).Return(nil)

	// act
	err := useCase.Handle("123", "456", 100, idempotencyKey.String())

	// assert
	assert.Error(t, err)
	assert.Equal(t, "from account not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestTransferAccountUseCase_Handle_ErrorGettingToAccount(t *testing.T) {
	// arrange
	mockRepo := new(usecases_mock.MockAccountRepository)
	mockBroker := new(usecases_mock.MockBroker)
	mockIdempotencyRepository := new(usecases_mock.MockIdempotencyRepository)

	useCase := NewTransferAccountUseCase(mockRepo, mockBroker, mockIdempotencyRepository)

	fromAcc := domain.NewAccount("123", "01234567890", "John Doe")
	mockRepo.On("GetAccountByNumber", "123").Return(fromAcc, nil)
	mockRepo.On("GetAccountByNumber", "456").Return((*domain.Account)(nil), errors.New("generic error"))

	mockBroker.On("Produce", mock.Anything, mock.Anything).Return(nil)

	idempotencyKey, _ := uuid.NewUUID()

	mockIdempotencyRepository.On("HasKey", idempotencyKey.String()).Return(false, nil)
	mockIdempotencyRepository.On("CreateKey", idempotencyKey.String()).Return(nil)

	// act
	err := useCase.Handle("123", "456", 100, idempotencyKey.String())

	// assert
	assert.Error(t, err)
	assert.Equal(t, "generic error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestTransferAccountUseCase_Handle_ToAccountNotFound(t *testing.T) {
	// arrange
	mockRepo := new(usecases_mock.MockAccountRepository)
	mockBroker := new(usecases_mock.MockBroker)
	mockIdempotencyRepository := new(usecases_mock.MockIdempotencyRepository)

	useCase := NewTransferAccountUseCase(mockRepo, mockBroker, mockIdempotencyRepository)

	fromAcc := domain.NewAccount("123", "01234567890", "John Doe")
	mockRepo.On("GetAccountByNumber", "123").Return(fromAcc, nil)
	mockRepo.On("GetAccountByNumber", "456").Return((*domain.Account)(nil), nil)

	mockBroker.On("Produce", mock.Anything, mock.Anything).Return(nil)

	idempotencyKey, _ := uuid.NewUUID()

	mockIdempotencyRepository.On("HasKey", idempotencyKey.String()).Return(false, nil)
	mockIdempotencyRepository.On("CreateKey", idempotencyKey.String()).Return(nil)

	// act
	err := useCase.Handle("123", "456", 100, idempotencyKey.String())

	// assert
	assert.Error(t, err)
	assert.Equal(t, "to account not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestTransferAccountUseCase_Handle_ErrorUpdatingToAccountBalance(t *testing.T) {
	// arrange
	mockRepo := new(usecases_mock.MockAccountRepository)
	mockBroker := new(usecases_mock.MockBroker)
	mockIdempotencyRepository := new(usecases_mock.MockIdempotencyRepository)
	useCase := NewTransferAccountUseCase(mockRepo, mockBroker, mockIdempotencyRepository)

	fromAcc := domain.NewAccount("123", "01234567890", "John Doe")
	toAcc := domain.NewAccount("456", "09876543210", "Jane Doe")
	mockRepo.On("GetAccountByNumber", "123").Return(fromAcc, nil)
	mockRepo.On("GetAccountByNumber", "456").Return(toAcc, nil)

	mockRepo.On("UpdateAccountBalance", toAcc).Return(errors.New("update error"))

	mockBroker.On("Produce", mock.Anything, mock.Anything).Return(nil)

	idempotencyKey, _ := uuid.NewUUID()

	mockIdempotencyRepository.On("HasKey", idempotencyKey.String()).Return(false, nil)
	mockIdempotencyRepository.On("CreateKey", idempotencyKey.String()).Return(nil)

	// act
	err := useCase.Handle("123", "456", 100, idempotencyKey.String())

	// assert
	assert.Error(t, err)
	assert.Equal(t, "update error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestTransferAccountUseCase_Handle_ErrorUpdatingFromAccountBalance(t *testing.T) {
	// arrange
	mockRepo := new(usecases_mock.MockAccountRepository)
	mockBroker := new(usecases_mock.MockBroker)
	mockIdempotencyRepository := new(usecases_mock.MockIdempotencyRepository)

	useCase := NewTransferAccountUseCase(mockRepo, mockBroker, mockIdempotencyRepository)

	fromAcc := domain.NewAccount("123", "01234567890", "John Doe")
	toAcc := domain.NewAccount("456", "09876543210", "Jane Doe")
	mockRepo.On("GetAccountByNumber", "123").Return(fromAcc, nil)
	mockRepo.On("GetAccountByNumber", "456").Return(toAcc, nil)

	mockRepo.On("UpdateAccountBalance", toAcc).Return(nil)
	mockRepo.On("UpdateAccountBalance", fromAcc).Return(errors.New("update error"))

	mockBroker.On("Produce", mock.Anything, mock.Anything).Return(nil)

	idempotencyKey, _ := uuid.NewUUID()

	mockIdempotencyRepository.On("HasKey", idempotencyKey.String()).Return(false, nil)
	mockIdempotencyRepository.On("CreateKey", idempotencyKey.String()).Return(nil)

	// act
	err := useCase.Handle("123", "456", 100, idempotencyKey.String())

	// assert
	assert.Error(t, err)
	assert.Equal(t, "update error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestTransferAccountUseCase_Handle_Success(t *testing.T) {
	// arrange
	mockRepo := new(usecases_mock.MockAccountRepository)
	mockBroker := new(usecases_mock.MockBroker)
	mockIdempotencyRepository := new(usecases_mock.MockIdempotencyRepository)

	useCase := NewTransferAccountUseCase(mockRepo, mockBroker, mockIdempotencyRepository)

	fromAcc := domain.NewAccount("123", "01234567890", "John Doe")
	toAcc := domain.NewAccount("456", "09876543210", "Jane Doe")
	mockRepo.On("GetAccountByNumber", "123").Return(fromAcc, nil)
	mockRepo.On("GetAccountByNumber", "456").Return(toAcc, nil)

	mockRepo.On("UpdateAccountBalance", toAcc).Return(nil)
	mockRepo.On("UpdateAccountBalance", fromAcc).Return(nil)

	mockBroker.On("Produce", mock.Anything, mock.Anything).Return(nil)

	idempotencyKey, _ := uuid.NewUUID()

	mockIdempotencyRepository.On("HasKey", idempotencyKey.String()).Return(false, nil)
	mockIdempotencyRepository.On("CreateKey", idempotencyKey.String()).Return(nil)

	// act
	err := useCase.Handle("123", "456", 100, idempotencyKey.String())

	// assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTransferAccountUseCase_Handle_HasKeyAlreadyUsed(t *testing.T) {
	// arrange
	mockRepo := new(usecases_mock.MockAccountRepository)
	mockBroker := new(usecases_mock.MockBroker)
	mockIdempotencyRepository := new(usecases_mock.MockIdempotencyRepository)

	useCase := NewTransferAccountUseCase(mockRepo, mockBroker, mockIdempotencyRepository)

	idempotencyKey, _ := uuid.NewUUID()

	mockIdempotencyRepository.On("HasKey", idempotencyKey.String()).Return(true, nil)

	// act
	err := useCase.Handle("123", "456", 100, idempotencyKey.String())

	// assert
	assert.Error(t, err)
	assert.Equal(t, err, errors.New("idempotency key already processed"))

	mockRepo.AssertExpectations(t)
	mockBroker.AssertExpectations(t)
	mockIdempotencyRepository.AssertExpectations(t)
}

func TestTransferAccountUseCase_Handle_HasKeyGenericError(t *testing.T) {
	// arrange
	mockRepo := new(usecases_mock.MockAccountRepository)
	mockBroker := new(usecases_mock.MockBroker)
	mockIdempotencyRepository := new(usecases_mock.MockIdempotencyRepository)

	useCase := NewTransferAccountUseCase(mockRepo, mockBroker, mockIdempotencyRepository)

	idempotencyKey, _ := uuid.NewUUID()

	genericError := errors.New("generic error")
	mockIdempotencyRepository.On("HasKey", idempotencyKey.String()).Return(false, genericError)

	// act
	err := useCase.Handle("123", "456", 100, idempotencyKey.String())

	// assert
	assert.Error(t, err)
	assert.Equal(t, err, genericError)

	mockRepo.AssertExpectations(t)
	mockBroker.AssertExpectations(t)
	mockIdempotencyRepository.AssertExpectations(t)
}

func TestTransferAccountUseCase_Handle_SuccessErrorSavingIdempotencyKeyUsed(t *testing.T) {
	// arrange
	mockRepo := new(usecases_mock.MockAccountRepository)
	mockBroker := new(usecases_mock.MockBroker)
	mockIdempotencyRepository := new(usecases_mock.MockIdempotencyRepository)

	useCase := NewTransferAccountUseCase(mockRepo, mockBroker, mockIdempotencyRepository)

	fromAcc := domain.NewAccount("123", "01234567890", "John Doe")
	toAcc := domain.NewAccount("456", "09876543210", "Jane Doe")
	mockRepo.On("GetAccountByNumber", "123").Return(fromAcc, nil)
	mockRepo.On("GetAccountByNumber", "456").Return(toAcc, nil)

	mockRepo.On("UpdateAccountBalance", toAcc).Return(nil)
	mockRepo.On("UpdateAccountBalance", fromAcc).Return(nil)

	mockBroker.On("Produce", mock.Anything, mock.Anything).Return(nil)

	idempotencyKey, _ := uuid.NewUUID()

	mockIdempotencyRepository.On("HasKey", idempotencyKey.String()).Return(false, nil)

	errorSavingUsedIdempotencyKey := errors.New("error saving used idempotency key")
	mockIdempotencyRepository.On("CreateKey", idempotencyKey.String()).Return(errorSavingUsedIdempotencyKey)

	// act
	err := useCase.Handle("123", "456", 100, idempotencyKey.String())

	// assert
	assert.Error(t, err)
	assert.Equal(t, err, errorSavingUsedIdempotencyKey)

	mockRepo.AssertExpectations(t)
	mockBroker.AssertExpectations(t)
	mockIdempotencyRepository.AssertExpectations(t)
}
