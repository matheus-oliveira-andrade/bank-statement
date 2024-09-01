package repositories

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/domain"
	"github.com/stretchr/testify/assert"
)

func getTestMovement() *domain.Movement {
	return &domain.Movement{
		Type:            "credit",
		AccountNumber:   "123456789",
		Value:           100.0,
		ToAccountNumber: "987654321",
		CreatedAt:       time.Now(),
	}
}

func TestCreateMovement_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMovementRepository(db)
	testMovement := getTestMovement()

	mock.ExpectExec("INSERT INTO movements").
		WithArgs(testMovement.Type, testMovement.AccountNumber, testMovement.Value, testMovement.ToAccountNumber, testMovement.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err = repo.CreateMovement(testMovement)

	// Assert
	assert.Nil(t, err)
}

func TestCreateMovement_DBError(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMovementRepository(db)
	testMovement := getTestMovement()

	mock.ExpectExec("INSERT INTO movements").
		WithArgs(testMovement.Type, testMovement.AccountNumber, testMovement.Value, testMovement.ToAccountNumber, testMovement.CreatedAt).
		WillReturnError(sql.ErrConnDone)

	// Act
	err = repo.CreateMovement(testMovement)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, sql.ErrConnDone, err)
}

func TestCreateMovement_NoRowsAffected(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMovementRepository(db)
	testMovement := getTestMovement()

	mock.ExpectExec("INSERT INTO movements").
		WithArgs(testMovement.Type, testMovement.AccountNumber, testMovement.Value, testMovement.ToAccountNumber, testMovement.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 0))

	// Act
	err = repo.CreateMovement(testMovement)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestCreateMovement_RowsAffectedError(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMovementRepository(db)
	testMovement := getTestMovement()

	mock.ExpectExec("INSERT INTO movements").
		WithArgs(testMovement.Type, testMovement.AccountNumber, testMovement.Value, testMovement.ToAccountNumber, testMovement.CreatedAt).
		WillReturnResult(sqlmock.NewErrorResult(sql.ErrConnDone))

	// Act
	err = repo.CreateMovement(testMovement)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, sql.ErrConnDone, err)
}
