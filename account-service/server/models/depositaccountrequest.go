package models

type DepositAccountRequest struct {
	Number         string `uri:"number" binding:"required"`
	Value          int64  `json:"value" binding:"required"`
	IdempotencyKey string `json:"idempotencyKey" binding:"required"`
}
