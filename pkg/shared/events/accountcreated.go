package events

type AccountCreated struct {
	Number   string `json:"number"`
	Name     string `json:"name"`
	Document string `json:"document"`
}

func NewAccountCreated(number, name, document string) *AccountCreated {
	return &AccountCreated{
		Number:   number,
		Name:     name,
		Document: document,
	}
}
