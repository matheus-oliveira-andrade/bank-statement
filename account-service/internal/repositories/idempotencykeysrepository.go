package repositories

import (
	"database/sql"
)

type IdempotencyKeysRepositoryInterface interface {
	HasKey(key string) (bool, error)
	CreateKey(key string) error
}

type IdempotencyKeysRepository struct {
	db *sql.DB
}

func NewIdempotencyKeysRepository(db *sql.DB) *IdempotencyKeysRepository {
	return &IdempotencyKeysRepository{
		db: db,
	}
}

func (r *IdempotencyKeysRepository) HasKey(key string) (bool, error) {
	tmp := ""
	err := r.db.QueryRow(`SELECT 1 FROM idempotencykeys WHERE Key = $1`, key).Scan(&tmp)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (r *IdempotencyKeysRepository) CreateKey(key string) error {
	result, err := r.db.Exec(`
	INSERT INTO idempotencykeys (Key, CreatedAt) VALUES ($1, CURRENT_TIMESTAMP)`, key)

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
