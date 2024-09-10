package eventhandlers_test

import (
	"errors"
	"testing"

	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/domain"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/eventhandlers"
	handlersmock "github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/eventhandlers/mocks"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/shared/events"
	"github.com/stretchr/testify/mock"
)

func TestStatementGenerationRequestedHandler_Handle_Success(t *testing.T) {
	// Arrange
	accountRepoMock := new(handlersmock.MockAccountRepository)
	statementGenRepoMock := new(handlersmock.MockStatementGenerationRepository)
	movementRepoMock := new(handlersmock.MockMovementRepository)
	documentGenApiMock := new(handlersmock.MockGenerateDocumentApi)
	templateCompilerMock := new(handlersmock.MockTemplateCompiler)

	handler := eventhandlers.NewStatementGenerationRequestedHandler(
		accountRepoMock, statementGenRepoMock, movementRepoMock, documentGenApiMock, templateCompilerMock,
	)

	account := &domain.Account{
		Document: "12345678900",
		Name:     "John Doe",
	}
	movements := []domain.Movement{
		{Type: string(domain.In), Value: 10000},
	}
	statementGeneration := &domain.StatementGeneration{}

	accountRepoMock.On("GetAccountByNumber", mock.Anything).Return(account, nil)
	statementGenRepoMock.On("GetStatementGeneration", mock.Anything).Return(statementGeneration, nil)
	movementRepoMock.On("GetMovements", mock.Anything).Return(&movements, nil)
	documentGenApiMock.On("GenerateFromHtml", mock.Anything).Return("pdf-data", nil)
	statementGenRepoMock.On("UpdateStatementGeneration", mock.Anything).Return(nil)
	templateCompilerMock.On("Compile", mock.Anything).Return("123XPTO321", nil)

	event := events.StatementGenerationRequested{
		AccountNumber: "12345678900",
	}

	// Act
	handler.Handle(event)

	// Assert
	accountRepoMock.AssertExpectations(t)
	statementGenRepoMock.AssertExpectations(t)
	movementRepoMock.AssertExpectations(t)
	documentGenApiMock.AssertExpectations(t)
}

func TestStatementGenerationRequestedHandler_Handle_ErrorOnGetStatementGeneration(t *testing.T) {
	// Arrange
	accountRepoMock := new(handlersmock.MockAccountRepository)
	statementGenRepoMock := new(handlersmock.MockStatementGenerationRepository)
	movementRepoMock := new(handlersmock.MockMovementRepository)
	documentGenApiMock := new(handlersmock.MockGenerateDocumentApi)
	templateCompilerMock := new(handlersmock.MockTemplateCompiler)

	handler := eventhandlers.NewStatementGenerationRequestedHandler(
		accountRepoMock, statementGenRepoMock, movementRepoMock, documentGenApiMock, templateCompilerMock,
	)

	statementGenRepoMock.On("GetStatementGeneration", mock.Anything).Return((*domain.StatementGeneration)(nil), errors.New("db error"))

	event := events.StatementGenerationRequested{
		AccountNumber: "12345678900",
	}

	// Act
	handler.Handle(event)

	// Assert
	statementGenRepoMock.AssertExpectations(t)
}

func TestStatementGenerationRequestedHandler_Handle_ErrorOnGetAccount(t *testing.T) {
	// Arrange
	accountRepoMock := new(handlersmock.MockAccountRepository)
	statementGenRepoMock := new(handlersmock.MockStatementGenerationRepository)
	movementRepoMock := new(handlersmock.MockMovementRepository)
	documentGenApiMock := new(handlersmock.MockGenerateDocumentApi)
	templateCompilerMock := new(handlersmock.MockTemplateCompiler)

	handler := eventhandlers.NewStatementGenerationRequestedHandler(
		accountRepoMock, statementGenRepoMock, movementRepoMock, documentGenApiMock, templateCompilerMock,
	)

	statementGeneration := &domain.StatementGeneration{}
	statementGenRepoMock.On("GetStatementGeneration", mock.Anything).Return(statementGeneration, nil)
	accountRepoMock.On("GetAccountByNumber", mock.Anything).Return((*domain.Account)(nil), errors.New("account not found"))
	statementGenRepoMock.On("UpdateStatementGeneration", mock.Anything).Return(nil)

	event := events.StatementGenerationRequested{
		AccountNumber: "12345678900",
	}

	// Act
	handler.Handle(event)

	// Assert
	accountRepoMock.AssertExpectations(t)
	statementGenRepoMock.AssertExpectations(t)
}

func TestStatementGenerationRequestedHandler_Handle_AccountNotFound(t *testing.T) {
	// Arrange
	accountRepoMock := new(handlersmock.MockAccountRepository)
	statementGenRepoMock := new(handlersmock.MockStatementGenerationRepository)
	movementRepoMock := new(handlersmock.MockMovementRepository)
	documentGenApiMock := new(handlersmock.MockGenerateDocumentApi)
	templateCompilerMock := new(handlersmock.MockTemplateCompiler)

	handler := eventhandlers.NewStatementGenerationRequestedHandler(
		accountRepoMock, statementGenRepoMock, movementRepoMock, documentGenApiMock, templateCompilerMock,
	)

	statementGeneration := &domain.StatementGeneration{}
	statementGenRepoMock.On("GetStatementGeneration", mock.Anything).Return(statementGeneration, nil)
	accountRepoMock.On("GetAccountByNumber", mock.Anything).Return((*domain.Account)(nil), nil)

	event := events.StatementGenerationRequested{
		AccountNumber: "12345678900",
	}

	// Act
	handler.Handle(event)

	// Assert
	accountRepoMock.AssertExpectations(t)
	statementGenRepoMock.AssertExpectations(t)
}

func TestStatementGenerationRequestedHandler_Handle_ErrorOnGetMovements(t *testing.T) {
	// Arrange
	accountRepoMock := new(handlersmock.MockAccountRepository)
	statementGenRepoMock := new(handlersmock.MockStatementGenerationRepository)
	movementRepoMock := new(handlersmock.MockMovementRepository)
	documentGenApiMock := new(handlersmock.MockGenerateDocumentApi)
	templateCompilerMock := new(handlersmock.MockTemplateCompiler)

	handler := eventhandlers.NewStatementGenerationRequestedHandler(
		accountRepoMock, statementGenRepoMock, movementRepoMock, documentGenApiMock, templateCompilerMock,
	)

	account := &domain.Account{}
	statementGeneration := &domain.StatementGeneration{}
	accountRepoMock.On("GetAccountByNumber", mock.Anything).Return(account, nil)
	statementGenRepoMock.On("GetStatementGeneration", mock.Anything).Return(statementGeneration, nil)
	movementRepoMock.On("GetMovements", mock.Anything).Return((*[]domain.Movement)(nil), errors.New("db error"))
	statementGenRepoMock.On("UpdateStatementGeneration", mock.Anything).Return(nil)

	event := events.StatementGenerationRequested{
		AccountNumber: "12345678900",
	}

	// Act
	handler.Handle(event)

	// Assert
	accountRepoMock.AssertExpectations(t)
	movementRepoMock.AssertExpectations(t)
}

func TestStatementGenerationRequestedHandler_Handle_MovementsNotFound(t *testing.T) {
	// Arrange
	accountRepoMock := new(handlersmock.MockAccountRepository)
	statementGenRepoMock := new(handlersmock.MockStatementGenerationRepository)
	movementRepoMock := new(handlersmock.MockMovementRepository)
	documentGenApiMock := new(handlersmock.MockGenerateDocumentApi)
	templateCompilerMock := new(handlersmock.MockTemplateCompiler)

	handler := eventhandlers.NewStatementGenerationRequestedHandler(
		accountRepoMock, statementGenRepoMock, movementRepoMock, documentGenApiMock, templateCompilerMock,
	)

	account := &domain.Account{}
	statementGeneration := &domain.StatementGeneration{}
	accountRepoMock.On("GetAccountByNumber", mock.Anything).Return(account, nil)
	statementGenRepoMock.On("GetStatementGeneration", mock.Anything).Return(statementGeneration, nil)
	movementRepoMock.On("GetMovements", mock.Anything).Return((*[]domain.Movement)(nil), nil)

	event := events.StatementGenerationRequested{
		AccountNumber: "12345678900",
	}

	// Act
	handler.Handle(event)

	// Assert
	accountRepoMock.AssertExpectations(t)
	movementRepoMock.AssertExpectations(t)
}

func TestStatementGenerationRequestedHandler_Handle_ErrorOnGenerateDocument(t *testing.T) {
	// Arrange
	accountRepoMock := new(handlersmock.MockAccountRepository)
	statementGenRepoMock := new(handlersmock.MockStatementGenerationRepository)
	movementRepoMock := new(handlersmock.MockMovementRepository)
	documentGenApiMock := new(handlersmock.MockGenerateDocumentApi)
	templateCompilerMock := new(handlersmock.MockTemplateCompiler)

	handler := eventhandlers.NewStatementGenerationRequestedHandler(
		accountRepoMock, statementGenRepoMock, movementRepoMock, documentGenApiMock, templateCompilerMock,
	)

	account := &domain.Account{}
	movements := []domain.Movement{}
	statementGeneration := &domain.StatementGeneration{}
	accountRepoMock.On("GetAccountByNumber", mock.Anything).Return(account, nil)
	statementGenRepoMock.On("GetStatementGeneration", mock.Anything).Return(statementGeneration, nil)
	movementRepoMock.On("GetMovements", mock.Anything).Return(&movements, nil)
	templateCompilerMock.On("Compile", mock.Anything).Return("123XPTO321", nil)
	documentGenApiMock.On("GenerateFromHtml", mock.Anything).Return("", errors.New("generation error"))
	statementGenRepoMock.On("UpdateStatementGeneration", mock.Anything).Return(nil)

	event := events.StatementGenerationRequested{
		AccountNumber: "12345678900",
	}

	// Act
	handler.Handle(event)

	// Assert
	documentGenApiMock.AssertExpectations(t)
}

func TestStatementGenerationRequestedHandler_Handle_ErrorOnUpdateStatementGeneration(t *testing.T) {
	// Arrange
	accountRepoMock := new(handlersmock.MockAccountRepository)
	statementGenRepoMock := new(handlersmock.MockStatementGenerationRepository)
	movementRepoMock := new(handlersmock.MockMovementRepository)
	documentGenApiMock := new(handlersmock.MockGenerateDocumentApi)
	templateCompilerMock := new(handlersmock.MockTemplateCompiler)

	handler := eventhandlers.NewStatementGenerationRequestedHandler(
		accountRepoMock, statementGenRepoMock, movementRepoMock, documentGenApiMock, templateCompilerMock,
	)

	account := &domain.Account{}
	movements := []domain.Movement{}
	statementGeneration := &domain.StatementGeneration{}
	accountRepoMock.On("GetAccountByNumber", mock.Anything).Return(account, nil)
	statementGenRepoMock.On("GetStatementGeneration", mock.Anything).Return(statementGeneration, nil)
	movementRepoMock.On("GetMovements", mock.Anything).Return(&movements, nil)
	documentGenApiMock.On("GenerateFromHtml", mock.Anything).Return("pdf-data", nil)
	templateCompilerMock.On("Compile", mock.Anything).Return("123XPTO321", nil)
	statementGenRepoMock.On("UpdateStatementGeneration", mock.Anything).Return(errors.New("update error"))

	event := events.StatementGenerationRequested{
		AccountNumber: "12345678900",
	}

	// Act
	handler.Handle(event)

	// Assert
	statementGenRepoMock.AssertExpectations(t)
}
