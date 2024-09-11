package repositories

import (
	"database/sql"

	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/internal/domain"
	"github.com/pkg/errors"
)

type StatementGenerationRepositoryInterface interface {
	CreateStatementGeneration(statementGeneration *domain.StatementGeneration) (string, error)
	HasStatementGenerationRunning(accountNumber string) (bool, error)
	GetStatementGeneration(accountNumber string) (*domain.StatementGeneration, error)
	UpdateStatementGeneration(statementGeneration *domain.StatementGeneration) error
	GetStatementGenerationById(id string) (*domain.StatementGeneration, error)
}

type StatementGenerationRepository struct {
	db *sql.DB
}

func NewStatementGenerationRepository(db *sql.DB) *StatementGenerationRepository {
	return &StatementGenerationRepository{
		db: db,
	}
}

func (r *StatementGenerationRepository) CreateStatementGeneration(statementGeneration *domain.StatementGeneration) (string, error) {
	row := r.db.QueryRow(`
	INSERT INTO statementsgeneration (Status, AccountNumber, CreatedAt, FinishedAt, Error, DocumentContent)
	VALUES ($1, $2, $3, $4, $5, $6)
	
	RETURNING Id
	`, statementGeneration.Status, statementGeneration.AccountNumber, statementGeneration.CreatedAt, statementGeneration.FinishedAt, statementGeneration.Error, statementGeneration.DocumentContent)

	var id string
	err := row.Scan(&id)

	return id, err
}

func (r *StatementGenerationRepository) HasStatementGenerationRunning(accountNumber string) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM statementsgeneration WHERE AccountNumber = $1 AND Status = $2)`

	err := r.db.QueryRow(query, accountNumber, domain.StatementGenerationRunnning).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (repo *StatementGenerationRepository) GetStatementGeneration(accountNumber string) (*domain.StatementGeneration, error) {
	query := `SELECT AccountNumber, Status, CreatedAt, FinishedAt, Error, DocumentContent FROM statementsgeneration WHERE AccountNumber = $1 AND Status = $2`
	row := repo.db.QueryRow(query, accountNumber, domain.StatementGenerationRunnning)

	var sg domain.StatementGeneration
	err := row.Scan(&sg.AccountNumber, &sg.Status, &sg.CreatedAt, &sg.FinishedAt, &sg.Error, &sg.DocumentContent)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to scan statement generation")
	}

	return &sg, nil
}

func (repo *StatementGenerationRepository) UpdateStatementGeneration(statementGeneration *domain.StatementGeneration) error {
	query := `
        UPDATE statementsgeneration
        SET Status = $1, FinishedAt = $2, Error = $3, DocumentContent = $4
        WHERE AccountNumber = $5
    `

	_, err := repo.db.Exec(query,
		statementGeneration.Status,
		statementGeneration.FinishedAt,
		statementGeneration.Error,
		statementGeneration.DocumentContent,
		statementGeneration.AccountNumber,
	)

	if err != nil {
		return errors.Wrap(err, "failed to update statement generation")
	}

	return nil
}

func (repo *StatementGenerationRepository) GetStatementGenerationById(id string) (*domain.StatementGeneration, error) {
	query := `SELECT Id, AccountNumber, Status, CreatedAt, FinishedAt, Error, DocumentContent FROM statementsgeneration WHERE Id = $1`
	row := repo.db.QueryRow(query, id)

	var sg domain.StatementGeneration
	err := row.Scan(&sg.Id, &sg.AccountNumber, &sg.Status, &sg.CreatedAt, &sg.FinishedAt, &sg.Error, &sg.DocumentContent)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to scan statement generation")
	}

	return &sg, nil
}
