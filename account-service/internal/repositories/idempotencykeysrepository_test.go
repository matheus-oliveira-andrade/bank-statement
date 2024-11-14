package repositories_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/matheus-oliveira-andrade/bank-statement/account-service/internal/repositories"
	"github.com/stretchr/testify/assert"
)

func TestHasKey_Found(t *testing.T) {
	// Arrange
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repositories.NewIdempotencyKeysRepository(db)

	expectedKey, _ := uuid.NewRandom()

	rows := sqlmock.NewRows([]string{"tmp"}).AddRow(1)
	mock.ExpectQuery("SELECT 1 FROM idempotencykeys WHERE Key = \\$1").
		WithArgs(expectedKey.String()).
		WillReturnRows(rows)

	// Act
	result, err := repo.HasKey(expectedKey.String())

	// Assert
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestHasKey_NotFound(t *testing.T) {
	// Arrange
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repositories.NewIdempotencyKeysRepository(db)

	expectedKey, _ := uuid.NewRandom()

	mock.ExpectQuery("SELECT 1 FROM idempotencykeys WHERE Key = \\$1").
		WithArgs(expectedKey.String()).
		WillReturnError(sql.ErrNoRows)

	// Act
	result, err := repo.HasKey(expectedKey.String())

	// Assert
	assert.NoError(t, err)
	assert.False(t, result)
}

func TestCreateKey_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewIdempotencyKeysRepository(db)

	expectedKey, _ := uuid.NewRandom()
	mock.ExpectExec("INSERT INTO idempotencykeys").
		WithArgs(expectedKey.String()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err = repo.CreateKey(expectedKey.String())

	// Assert
	assert.Nil(t, err)
}

func TestCreateKey_DBError(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewIdempotencyKeysRepository(db)

	expectedKey, _ := uuid.NewRandom()
	mock.ExpectExec("INSERT INTO idempotencykeys").
		WithArgs(expectedKey.String()).
		WillReturnError(sql.ErrConnDone)

	// Act
	err = repo.CreateKey(expectedKey.String())

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, sql.ErrConnDone, err)
}

func TestCreateKey_NoRowsAffected(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewIdempotencyKeysRepository(db)

	expectedKey, _ := uuid.NewRandom()
	mock.ExpectExec("INSERT INTO idempotencykeys").
		WithArgs(expectedKey.String()).
		WillReturnResult(sqlmock.NewResult(1, 0))

	// Act
	err = repo.CreateKey(expectedKey.String())

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
}
