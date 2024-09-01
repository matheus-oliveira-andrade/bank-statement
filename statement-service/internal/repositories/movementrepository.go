package repositories

import (
	"database/sql"

	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/domain"
)

type MovementRepositoryInterface interface {
	CreateMovement(movement *domain.Movement) error
}

type MovementRepository struct {
	db *sql.DB
}

func NewMovementRepository(db *sql.DB) *MovementRepository {
	return &MovementRepository{
		db: db,
	}
}

func (r *MovementRepository) CreateMovement(movement *domain.Movement) error {
	result, err := r.db.Exec(`
	INSERT INTO movements (Type, AccountNumber, Value, ToAccountNumber, CreatedAt)
	VALUES ($1, $2, $3, $4, $5)
	`, movement.Type, movement.AccountNumber, movement.Value, movement.ToAccountNumber, movement.CreatedAt)

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
