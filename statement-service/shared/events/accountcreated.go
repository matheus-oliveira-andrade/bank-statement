package events

type AccountCreated struct {
	Number   string `json:"number"`
	Name     string `json:"name"`
	Document string `json:"document"`
}
