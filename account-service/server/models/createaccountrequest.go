package models

type CreateAccountRequest struct {
	Document string `json:"document"`
	Name     string `json:"name"`
}
