package repositories

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestCreateStatementGeneration_Error(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewStatementGenerationRepository(db)

	statementGeneration := &domain.StatementGeneration{
		Status:        "Running",
		AccountNumber: "123456",
		CreatedAt:     time.Now(),
	}

	mock.ExpectQuery(`INSERT INTO statementsgeneration`).
		WithArgs(statementGeneration.Status, statementGeneration.AccountNumber, statementGeneration.CreatedAt).
		WillReturnError(sqlmock.ErrCancelled)

	// act
	id, err := repo.CreateStatementGeneration(statementGeneration)

	// assert
	assert.Error(t, err)
	assert.Empty(t, id)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateStatementGeneration_Success(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewStatementGenerationRepository(db)

	statementGeneration := &domain.StatementGeneration{
		Status:        "Running",
		AccountNumber: "123456",
		CreatedAt:     time.Now(),
	}

	mock.ExpectQuery(`INSERT INTO statementsgeneration`).
		WithArgs(statementGeneration.Status, statementGeneration.AccountNumber, statementGeneration.CreatedAt).
		WillReturnRows(sqlmock.NewRows([]string{"Id"}).AddRow("1"))

	// act
	id, err := repo.CreateStatementGeneration(statementGeneration)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, "1", id)
}

func TestHasStatementGenerationRunning_Found(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewStatementGenerationRepository(db)

	accountNumber := "123456"
	status := domain.StatementGenerationRunnning

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM statementsgeneration WHERE AccountNumber = \$1 AND Status = \$2\)`).
		WithArgs(accountNumber, status).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	// act
	exists, err := repo.HasStatementGenerationRunning(accountNumber)

	// assert
	assert.NoError(t, err)
	assert.True(t, exists)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestHasStatementGenerationRunning_NotFound(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewStatementGenerationRepository(db)

	accountNumber := "123456"
	status := domain.StatementGenerationRunnning

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM statementsgeneration WHERE AccountNumber = \$1 AND Status = \$2\)`).
		WithArgs(accountNumber, status).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}))

	// act
	exists, err := repo.HasStatementGenerationRunning(accountNumber)

	// assert
	assert.NoError(t, err)
	assert.False(t, exists)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetStatementGeneration_Found(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewStatementGenerationRepository(db)

	accountNumber := "123456"
	expectedSG := domain.StatementGeneration{
		AccountNumber:   accountNumber,
		Status:          "running",
		CreatedAt:       time.Now(),
		FinishedAt:      time.Now().Add(1 * time.Hour),
		Error:           "none",
		DocumentContent: "PDF content",
	}

	mock.ExpectQuery(`SELECT AccountNumber, Status, CreatedAt, FinishedAt, Error, DocumentContent FROM statementsgeneration WHERE AccountNumber = \$1`).
		WithArgs(accountNumber).
		WillReturnRows(sqlmock.NewRows([]string{"AccountNumber", "Status", "CreatedAt", "FinishedAt", "Error", "DocumentContent"}).
			AddRow(expectedSG.AccountNumber, expectedSG.Status, expectedSG.CreatedAt, expectedSG.FinishedAt, expectedSG.Error, expectedSG.DocumentContent))

	// act
	sg, err := repo.GetStatementGeneration(accountNumber)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, &expectedSG, sg)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetStatementGeneration_NotFound(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewStatementGenerationRepository(db)

	accountNumber := "123456"

	mock.ExpectQuery(`SELECT AccountNumber, Status, CreatedAt, FinishedAt, Error, DocumentContent FROM statementsgeneration WHERE AccountNumber = \$1`).
		WithArgs(accountNumber).
		WillReturnRows(sqlmock.NewRows([]string{"AccountNumber", "Status", "CreatedAt", "FinishedAt", "Error", "DocumentContent"}))

	// act
	sg, err := repo.GetStatementGeneration(accountNumber)

	// assert
	assert.NoError(t, err)
	assert.Nil(t, sg) // Nenhum resultado encontrado

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetStatementGeneration_Error(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewStatementGenerationRepository(db)

	accountNumber := "123456"

	mock.ExpectQuery(`SELECT AccountNumber, Status, CreatedAt, FinishedAt, Error, DocumentContent FROM statementsgeneration WHERE AccountNumber = \$1`).
		WithArgs(accountNumber).
		WillReturnError(errors.New("query error"))

	// act
	sg, err := repo.GetStatementGeneration(accountNumber)

	// assert
	assert.Error(t, err)
	assert.Nil(t, sg)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateStatementGeneration_Success(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewStatementGenerationRepository(db)

	statementGeneration := &domain.StatementGeneration{
		AccountNumber:   "123456",
		Status:          "completed",
		CreatedAt:       time.Now(),
		FinishedAt:      time.Now().Add(1 * time.Hour),
		Error:           "none",
		DocumentContent: "Updated PDF content",
	}

	mock.ExpectExec(`UPDATE statementsgeneration SET Status = \$1, FinishedAt = \$2, Error = \$3, DocumentContent = \$4 WHERE AccountNumber = \$5`).
		WithArgs(
			statementGeneration.Status,
			statementGeneration.FinishedAt,
			statementGeneration.Error,
			statementGeneration.DocumentContent,
			statementGeneration.AccountNumber,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// act
	err = repo.UpdateStatementGeneration(statementGeneration)

	// assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateStatementGeneration_Error(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewStatementGenerationRepository(db)

	statementGeneration := &domain.StatementGeneration{
		AccountNumber:   "123456",
		Status:          "completed",
		CreatedAt:       time.Now(),
		FinishedAt:      time.Now().Add(1 * time.Hour),
		Error:           "none",
		DocumentContent: "Updated PDF content",
	}

	mock.ExpectExec(`UPDATE statementsgeneration SET Status = \$1, FinishedAt = \$2, Error = \$3, DocumentContent = \$4 WHERE AccountNumber = \$5`).
		WithArgs(
			statementGeneration.Status,
			statementGeneration.FinishedAt,
			statementGeneration.Error,
			statementGeneration.DocumentContent,
			statementGeneration.AccountNumber,
		).
		WillReturnError(errors.New("update error"))

	// act
	err = repo.UpdateStatementGeneration(statementGeneration)

	// assert
	assert.Error(t, err)
	assert.EqualError(t, err, "failed to update statement generation: update error")

	assert.NoError(t, mock.ExpectationsWereMet())
}
