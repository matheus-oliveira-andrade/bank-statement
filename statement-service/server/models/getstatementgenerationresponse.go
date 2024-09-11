package models

type GetStatementGenerationResponse struct {
	File string `json:"file"`
}

func NewGetStatementGenerationResponse(file string) *GetStatementGenerationResponse {
	return &GetStatementGenerationResponse{
		File: file,
	}
}
