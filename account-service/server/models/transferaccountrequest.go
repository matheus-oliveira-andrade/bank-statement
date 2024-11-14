package models

type TransferAccountRequest struct {
	FromNumber     string `uri:"fromNumber" binding:"required"`
	ToNumber       string `uri:"toNumber" binding:"required"`
	Value          int64  `json:"value" binding:"required"`
	IdempotencyKey string `json:"idempotencyKey" binding:"required"`
}
