package repositories

import (
	"database/sql"

	"github.com/matheus-oliveira-andrade/bank-statement/account-service/internal/domain"
)

type AccountRepositoryInterface interface {
	GetAccountByNumber(number string) (*domain.Account, error)
	GetAccountByDocument(document string) (*domain.Account, error)
	GetNextAccountNumber() (string, error)
	CreateAccount(account *domain.Account) (string, error)
	UpdateAccountBalance(account *domain.Account) error
}

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (r *AccountRepository) GetAccountByNumber(number string) (*domain.Account, error) {
	row := r.db.QueryRow(`
		SELECT Id, Number, Name, Document, Balance, CreatedAt, UpdatedAt
		FROM accounts 
		WHERE Number = $1
	`, number)

	var account domain.Account
	err := row.Scan(&account.Id, &account.Number, &account.Name, &account.Document, &account.Balance, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &account, nil
}

func (r *AccountRepository) GetAccountByDocument(document string) (*domain.Account, error) {
	row := r.db.QueryRow(`
		SELECT Id, Number, Name, Document, Balance, CreatedAt, UpdatedAt
		FROM accounts 
		WHERE Document = $1
	`, document)

	var account domain.Account
	err := row.Scan(&account.Id, &account.Number, &account.Name, &account.Document, &account.Balance, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &account, nil
}

func (r *AccountRepository) GetNextAccountNumber() (string, error) {
	var number string

	row := r.db.QueryRow(`SELECT CAST(coalesce(MAX(Number), '0') AS integer) + 1 FROM accounts`)

	err := row.Scan(&number)

	return number, err
}

func (r *AccountRepository) CreateAccount(account *domain.Account) (string, error) {

	row := r.db.QueryRow(`
	INSERT INTO accounts (Number, Name, Document, Balance, CreatedAt, UpdatedAt)
	VALUES ($1, $2, $3, $4, $5, $6)
	
	RETURNING Id
	`, account.Number, account.Name, account.Document, account.Balance, account.CreatedAt, account.UpdatedAt)

	var id string
	err := row.Scan(&id)

	return id, err
}

func (r *AccountRepository) UpdateAccountBalance(account *domain.Account) error {
	result, err := r.db.Exec(`UPDATE accounts SET Balance = $1, UpdatedAt = $2 WHERE Id = $3`,
		account.Balance, account.UpdatedAt, account.Id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
