package usecases

import (
	"errors"
	"testing"

	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/domain"
	usecases_mocks "github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/usecases/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHandle_GetStatementGeneration_Success(t *testing.T) {
	// arrange
	statementRepositoryMock := new(usecases_mocks.MockStatementGenerationRepository)

	statementGeneration, _ := domain.NewStatementGeneration("123")
	statementGeneration.SetAsGenerated("dwqeqciaosada1")

	id := "12"

	statementRepositoryMock.On("GetStatementGenerationById", id).Return(statementGeneration, nil)

	usecase := NewGetStatementGenerationUseCase(statementRepositoryMock)

	// act
	result, err := usecase.Handle(id)

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, statementGeneration.DocumentContent, result)

	statementRepositoryMock.AssertExpectations(t)
}

func TestHandle_GetStatementGeneration_SuccessWithError(t *testing.T) {
	// arrange
	statementRepositoryMock := new(usecases_mocks.MockStatementGenerationRepository)

	statementGeneration, _ := domain.NewStatementGeneration("123")
	statementGeneration.SetAsGeneratedWithError(errors.New("fake error"))

	id := "12"

	statementRepositoryMock.On("GetStatementGenerationById", id).Return(statementGeneration, nil)

	usecase := NewGetStatementGenerationUseCase(statementRepositoryMock)

	// act
	result, err := usecase.Handle(id)

	// assert
	assert.Error(t, err)
	assert.Empty(t, result)

	statementRepositoryMock.AssertExpectations(t)
}

func TestHandle_GetStatementGeneration_Error(t *testing.T) {
	// arrange
	statementRepositoryMock := new(usecases_mocks.MockStatementGenerationRepository)

	id := "12"

	statementRepositoryMock.On("GetStatementGenerationById", id).Return((*domain.StatementGeneration)(nil), errors.New("fake error"))

	usecase := NewGetStatementGenerationUseCase(statementRepositoryMock)

	// act
	result, err := usecase.Handle(id)

	// assert
	assert.Error(t, err)
	assert.Empty(t, result)

	statementRepositoryMock.AssertExpectations(t)
}
