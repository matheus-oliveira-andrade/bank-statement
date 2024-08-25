package models

type GetAccountRequest struct {
	Number string `uri:"number" binding:"required"`
}
