package events

const StatementGenerationRequestedEventKey = "StatementGenerationRequested"

type StatementGenerationRequested struct {
	Id            string `json:"id"`
	AccountNumber string `json:"accountNumber"`
}

func NewStatementGenerationRequested(id string, accountNumber string) *StatementGenerationRequested {
	return &StatementGenerationRequested{
		Id:            id,
		AccountNumber: accountNumber,
	}
}
