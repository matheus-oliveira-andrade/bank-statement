package models

type TriggerStatementGenerationResponse struct {
	FileBase64 string `json:"fileBase64"`
}

func NewTriggerStatementGenerationResponse(fileBase64 string) *TriggerStatementGenerationResponse {
	return &TriggerStatementGenerationResponse{
		FileBase64: fileBase64,
	}
}
