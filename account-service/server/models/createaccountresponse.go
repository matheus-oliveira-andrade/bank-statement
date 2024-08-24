package models

type CreateAccountResponse struct {
	Number string `json:"number"`
}

func NewCreateAccountResponse(number string) *CreateAccountResponse {
	return &CreateAccountResponse{
		Number: number,
	}
}
