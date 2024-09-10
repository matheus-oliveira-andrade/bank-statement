package repositories

import (
	"database/sql"
	"errors"
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

func TestGetMovements_NotFound(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMovementRepository(db)

	accountNumber := "123456"

	mock.ExpectQuery(`SELECT Type, AccountNumber, Value, ToAccountNumber, CreatedAt FROM movements WHERE AccountNumber = \$1`).
		WithArgs(accountNumber).
		WillReturnRows(sqlmock.NewRows([]string{"Type", "AccountNumber", "Value", "ToAccountNumber", "CreatedAt"}))

	// act
	movements, err := repo.GetMovements(accountNumber)

	// assert
	assert.NoError(t, err)
	assert.Empty(t, movements)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetMovements_FoundMovements(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMovementRepository(db)

	accountNumber := "123456"

	rows := sqlmock.NewRows([]string{"Type", "AccountNumber", "Value", "ToAccountNumber", "CreatedAt"}).
		AddRow("deposit", "123456", 100.0, "654321", time.Now()).
		AddRow("withdrawal", "123456", 50.0, "654321", time.Now())

	mock.ExpectQuery(`SELECT Type, AccountNumber, Value, ToAccountNumber, CreatedAt FROM movements WHERE AccountNumber = \$1`).
		WithArgs(accountNumber).
		WillReturnRows(rows)

	// act
	movements, err := repo.GetMovements(accountNumber)

	// assert
	assert.NoError(t, err)
	assert.NotEmpty(t, movements)
	assert.Equal(t, 2, len(*movements))

	assert.Equal(t, "deposit", (*movements)[0].Type)
	assert.Equal(t, "withdrawal", (*movements)[1].Type)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetMovements_ErrorReadingMovements(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMovementRepository(db)

	accountNumber := "123456"

	mock.ExpectQuery(`SELECT Type, AccountNumber, Value, ToAccountNumber, CreatedAt FROM movements WHERE AccountNumber = \$1`).
		WithArgs(accountNumber).
		WillReturnError(errors.New("query failed"))

	// act
	movements, err := repo.GetMovements(accountNumber)

	// assert
	assert.Error(t, err)
	assert.Empty(t, movements)

	assert.NoError(t, mock.ExpectationsWereMet())
}
