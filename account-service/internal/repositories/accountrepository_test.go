package repositories

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/matheus-oliveira-andrade/bank-statement/account-service/internal/domain"
	"github.com/stretchr/testify/assert"
)

// helper para criar uma conta esperada
func getExpectedAccount() *domain.Account {
	return &domain.Account{
		Id:        "1",
		Number:    "123456789",
		Name:      "John Doe",
		Document:  "12345678901",
		Balance:   1000.0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestGetAccountByNumber_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewAccountRepository(db)
	expectedAccount := getExpectedAccount()

	rows := sqlmock.NewRows([]string{"Id", "Number", "Name", "Document", "Balance", "CreatedAt", "UpdatedAt"}).
		AddRow(expectedAccount.Id, expectedAccount.Number, expectedAccount.Name, expectedAccount.Document, expectedAccount.Balance, expectedAccount.CreatedAt, expectedAccount.UpdatedAt)
	mock.ExpectQuery("SELECT Id, Number, Name, Document, Balance, CreatedAt, UpdatedAt FROM accounts WHERE Number = \\$1").
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

	mock.ExpectQuery("SELECT Id, Number, Name, Document, Balance, CreatedAt, UpdatedAt FROM accounts WHERE Number = \\$1").
		WithArgs("987654321").
		WillReturnError(sql.ErrNoRows)

	account, err := repo.GetAccountByNumber("987654321")
	assert.Error(t, err, "an error was expected when account is not found")
	assert.Nil(t, account, "account should be nil when not found")
}

func TestGetAccountByNumber_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewAccountRepository(db)

	mock.ExpectQuery("SELECT Id, Number, Name, Document, Balance, CreatedAt, UpdatedAt FROM accounts WHERE Number = \\$1").
		WithArgs("123456789").
		WillReturnError(sql.ErrConnDone)

	account, err := repo.GetAccountByNumber("123456789")
	assert.Error(t, err, "an error was expected due to a database connection issue")
	assert.Nil(t, account, "account should be nil when there is a database error")
}
