package domain

import (
	"errors"
	"time"
)

const (
	StatementGenerationRunnning = "running"
	StatementGenerationFinished = "finished"
	StatementGenerationError    = "errorGenerating"
)

type StatementGeneration struct {
	Id              string
	Status          string
	AccountNumber   string
	CreatedAt       time.Time
	FinishedAt      time.Time
	Error           string
	DocumentContent string
}

func NewStatementGeneration(accountNumber string) (*StatementGeneration, error) {
	if accountNumber == "" {
		return nil, errors.New("account number is required")
	}

	return &StatementGeneration{
		AccountNumber: accountNumber,
		Status:        StatementGenerationRunnning,
		CreatedAt:     time.Now(),
	}, nil
}

func (sg *StatementGeneration) SetAsGenerated(report string) {
	sg.DocumentContent = report
	sg.Status = StatementGenerationFinished
	sg.FinishedAt = time.Now()
}

func (sg *StatementGeneration) SetAsGeneratedWithError(err error) {
	sg.Error = err.Error()
	sg.Status = StatementGenerationError
	sg.FinishedAt = time.Now()
}
