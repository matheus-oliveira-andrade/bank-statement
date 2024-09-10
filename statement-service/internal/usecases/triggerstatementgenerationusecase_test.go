package usecases

import (
	"testing"

	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/domain"
	usecases_mocks "github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/usecases/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandle_StatementGenerationRunning(t *testing.T) {
	// arrange
	mockStatementRepo := new(usecases_mocks.MockStatementGenerationRepository)
	mockAccountRepo := new(usecases_mocks.MockAccountRepository)
	mockBroker := new(usecases_mocks.MockBroker)

	useCase := NewTriggerStatementGenerationUseCase(mockStatementRepo, mockAccountRepo, mockBroker)

	acc := &domain.Account{
		Number: "123456",
	}

	mockAccountRepo.On("GetAccountByNumber", mock.Anything).Return(acc, nil)

	accountNumber := "123456"
	mockStatementRepo.On("HasStatementGenerationRunning", accountNumber).Return(true, nil)

	mockBroker.On("Produce", mock.Anything, mock.Anything).Return(nil)

	// act
	result, err := useCase.Handle(accountNumber)

	// assert
	assert.Error(t, err)
	assert.Equal(t, "already have a statement generation running", err.Error())
	assert.Equal(t, "", result)
	mockStatementRepo.AssertExpectations(t)
}

func TestHandle_ErrorCheckingGenerationRunning(t *testing.T) {
	// arrange
	mockStatementRepo := new(usecases_mocks.MockStatementGenerationRepository)
	mockAccountRepo := new(usecases_mocks.MockAccountRepository)
	mockBroker := new(usecases_mocks.MockBroker)

	useCase := NewTriggerStatementGenerationUseCase(mockStatementRepo, mockAccountRepo, mockBroker)

	acc := &domain.Account{
		Number: "123456",
	}

	mockAccountRepo.On("GetAccountByNumber", mock.Anything).Return(acc, nil)

	accountNumber := "123456"
	mockStatementRepo.On("HasStatementGenerationRunning", accountNumber).Return(false, assert.AnError)

	mockBroker.On("Produce", mock.Anything, mock.Anything).Return(nil)

	// act
	result, err := useCase.Handle(accountNumber)

	// assert
	assert.Error(t, err)
	assert.Equal(t, "error search for statement generation running", err.Error())
	assert.Equal(t, "", result)
	mockStatementRepo.AssertExpectations(t)
}

func TestHandle_ErrorCreatingStatementGeneration(t *testing.T) {
	// arrange
	mockStatementRepo := new(usecases_mocks.MockStatementGenerationRepository)
	mockAccountRepo := new(usecases_mocks.MockAccountRepository)
	mockBroker := new(usecases_mocks.MockBroker)

	useCase := NewTriggerStatementGenerationUseCase(mockStatementRepo, mockAccountRepo, mockBroker)

	acc := &domain.Account{
		Number: "123456",
	}

	mockAccountRepo.On("GetAccountByNumber", mock.Anything).Return(acc, nil)

	accountNumber := "123456"
	mockStatementRepo.On("HasStatementGenerationRunning", accountNumber).Return(false, nil)
	mockStatementRepo.On("CreateStatementGeneration", mock.Anything).Return("", assert.AnError)

	mockBroker.On("Produce", mock.Anything, mock.Anything).Return(nil)

	// act
	result, err := useCase.Handle(accountNumber)

	// assert
	assert.Error(t, err)
	assert.Equal(t, "error creating statement generation", err.Error())
	assert.Equal(t, "", result)
	mockStatementRepo.AssertExpectations(t)
}

func TestHandle_Success(t *testing.T) {
	// arrange
	mockStatementRepo := new(usecases_mocks.MockStatementGenerationRepository)
	mockAccountRepo := new(usecases_mocks.MockAccountRepository)
	mockBroker := new(usecases_mocks.MockBroker)

	useCase := NewTriggerStatementGenerationUseCase(mockStatementRepo, mockAccountRepo, mockBroker)

	acc := &domain.Account{
		Number: "123456",
	}

	mockAccountRepo.On("GetAccountByNumber", mock.Anything).Return(acc, nil)

	accountNumber := "123456"
	triggerId := "abc123"
	statementGeneration, _ := domain.NewStatementGeneration(accountNumber)
	mockStatementRepo.On("HasStatementGenerationRunning", accountNumber).Return(false, nil)
	mockStatementRepo.On("CreateStatementGeneration", statementGeneration).Return(triggerId, nil)

	mockBroker.On("Produce", mock.Anything, mock.Anything).Return(nil)

	// act
	result, err := useCase.Handle(accountNumber)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, triggerId, result)
	mockStatementRepo.AssertExpectations(t)
}

func TestHandle_AccountNotFound(t *testing.T) {
	// arrange
	mockStatementRepo := new(usecases_mocks.MockStatementGenerationRepository)
	mockAccountRepo := new(usecases_mocks.MockAccountRepository)
	mockBroker := new(usecases_mocks.MockBroker)

	useCase := NewTriggerStatementGenerationUseCase(mockStatementRepo, mockAccountRepo, mockBroker)

	mockAccountRepo.On("GetAccountByNumber", mock.Anything).Return((*domain.Account)(nil), nil)

	accountNumber := "123456"

	mockBroker.On("Produce", mock.Anything, mock.Anything).Return(nil)

	// act
	result, err := useCase.Handle(accountNumber)

	// assert
	assert.Error(t, err)
	assert.Equal(t, "account not found: 123456", err.Error())
	assert.Empty(t, result)

	mockStatementRepo.AssertExpectations(t)
}
