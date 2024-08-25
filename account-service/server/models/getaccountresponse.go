package models

import (
	"time"

	"github.com/matheus-oliveira-andrade/bank-statement/account-service/internal/domain"
)

type GetAccountResponse struct {
	Number    string    `json:"number"`
	Document  string    `json:"document"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewGetAccountResponse(acc *domain.Account) *GetAccountResponse {
	return &GetAccountResponse{
		Number:    acc.Number,
		Document:  acc.Document,
		Balance:   acc.Balance,
		CreatedAt: acc.CreatedAt,
		UpdatedAt: acc.UpdatedAt,
	}
}
