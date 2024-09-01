package domain

import "time"

type MovementType string

const (
	In  MovementType = "in"
	Out MovementType = "out"
)

type Movement struct {
	Id              int
	Type            string
	AccountNumber   string
	Value           int64
	ToAccountNumber string
	CreatedAt       time.Time
}

func NewDepositedFundsMovement(accountNumber string, value int64) *Movement {
	return &Movement{
		Type:          string(In),
		AccountNumber: accountNumber,
		Value:         value,
		CreatedAt:     time.Now(),
	}
}
