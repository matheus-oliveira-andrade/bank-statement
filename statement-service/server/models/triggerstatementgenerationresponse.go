package models

type TriggerStatementGenerationResponse struct {
	TriggerId string
}

func NewTriggerStatementGenerationResponse(triggerId string) *TriggerStatementGenerationResponse {
	return &TriggerStatementGenerationResponse{
		TriggerId: triggerId,
	}
}
