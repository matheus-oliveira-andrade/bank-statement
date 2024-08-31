package repositories

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/domain"
	"github.com/stretchr/testify/assert"
)

func getExpectedAccount() *domain.Account {
	return &domain.Account{
		Number:   "123456789",
		Name:     "John Doe",
		Document: "12345678901",
		Balance:  1000.0,
	}
}

func TestGetAccountByNumber_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewAccountRepository(db)
	expectedAccount := getExpectedAccount()

	rows := sqlmock.NewRows([]string{"Number", "Name", "Document", "Balance"}).
		AddRow(expectedAccount.Number, expectedAccount.Name, expectedAccount.Document, expectedAccount.Balance)
	mock.ExpectQuery("SELECT Number, Name, Document, Balance FROM accounts WHERE Number = \\$1").
		WithArgs(expectedAccount.Number).
		WillReturnRows(rows)

	account, err := repo.GetAccountByNumber(expectedAccount.Number)
	assert.NoError(t, err, "error was not expected while getting account by number")
	assert.NotNil(t, account, "account should not be nil")
	assert.Equal(t, expectedAccount, account)
}

func TestGetAccountByNumber_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewAccountRepository(db)

	mock.ExpectQuery("SELECT Number, Name, Document, Balance FROM accounts WHERE Number = \\$1").
		WithArgs("987654321").
		WillReturnError(sql.ErrNoRows)

	account, err := repo.GetAccountByNumber("987654321")
	assert.Nil(t, err)
	assert.Nil(t, account)
}

func TestGetAccountByNumber_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewAccountRepository(db)

	mock.ExpectQuery("SELECT Number, Name, Document, Balance FROM accounts WHERE Number = \\$1").
		WithArgs("123456789").
		WillReturnError(sql.ErrConnDone)

	account, err := repo.GetAccountByNumber("123456789")
	assert.Error(t, err, "an error was expected due to a database connection issue")
	assert.Nil(t, account, "account should be nil when there is a database error")
}

func TestUpdateAccountBalance_ErrOnExecute(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewAccountRepository(db)

	acc := domain.NewAccount("1", "12345678901", "John Dii")

	mock.ExpectExec("UPDATE accounts SET Balance = \\$1 WHERE Number = \\$2").
		WithArgs(acc.Balance, acc.Number).
		WillReturnError(sql.ErrConnDone)

	// Act
	err = repo.UpdateAccountBalance(acc)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, sql.ErrConnDone, err)
}

func TestUpdateAccountBalance_NotRowsAffected(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewAccountRepository(db)

	acc := domain.NewAccount("1", "12345678901", "John Dii")

	mock.ExpectExec("UPDATE accounts SET Balance = \\$1 WHERE Number = \\$2").
		WithArgs(acc.Balance, acc.Number).
		WillReturnResult(sqlmock.NewResult(1, 0))

	// Act
	err = repo.UpdateAccountBalance(acc)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestUpdateAccountBalance_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewAccountRepository(db)

	acc := domain.NewAccount("1", "12345678901", "John Dii")

	mock.ExpectExec("UPDATE accounts SET Balance = \\$1 WHERE Number = \\$2").
		WithArgs(acc.Balance, acc.Number).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err = repo.UpdateAccountBalance(acc)

	// Assert
	assert.Nil(t, err)
}
