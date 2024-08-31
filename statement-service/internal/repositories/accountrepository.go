package repositories

import (
	"database/sql"

	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/domain"
)

type AccountRepositoryInterface interface {
	GetAccountByNumber(number string) (*domain.Account, error)
	CreateAccount(account *domain.Account) error
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
		SELECT Number, Name, Document, Balance
		FROM accounts 
		WHERE Number = $1
	`, number)

	var account domain.Account
	err := row.Scan(&account.Number, &account.Name, &account.Document, &account.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &account, nil
}

func (r *AccountRepository) CreateAccount(account *domain.Account) error {
	result, err := r.db.Exec(`
	INSERT INTO accounts (Number, Name, Document, Balance)
	VALUES ($1, $2, $3, $4)
	`, account.Number, account.Name, account.Document, account.Balance)

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

func (r *AccountRepository) UpdateAccountBalance(account *domain.Account) error {
	result, err := r.db.Exec(`UPDATE accounts SET Balance = $1 WHERE Number = $2`,
		account.Balance, account.Number)

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
