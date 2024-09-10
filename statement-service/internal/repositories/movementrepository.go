package repositories

import (
	"database/sql"

	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/domain"
	"github.com/pkg/errors"
)

type MovementRepositoryInterface interface {
	CreateMovement(movement *domain.Movement) error
	GetMovements(accountNumber string) (*[]domain.Movement, error)
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

func (r *MovementRepository) GetMovements(accountNumber string) (*[]domain.Movement, error) {
	query := `SELECT Type, AccountNumber, Value, ToAccountNumber, CreatedAt FROM movements WHERE AccountNumber = $1`
	rows, err := r.db.Query(query, accountNumber)

	if err != nil {
		return nil, errors.Wrap(err, "failed to execute query")
	}
	defer rows.Close()

	var movements []domain.Movement

	for rows.Next() {
		var sg domain.Movement
		err := rows.Scan(&sg.Type, &sg.AccountNumber, &sg.Value, &sg.ToAccountNumber, &sg.CreatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan movement")
		}
		movements = append(movements, sg)
	}

	if err = rows.Err(); err != nil {
		return &[]domain.Movement{}, errors.Wrap(err, "error iterating over result rows")
	}

	return &movements, nil
}
