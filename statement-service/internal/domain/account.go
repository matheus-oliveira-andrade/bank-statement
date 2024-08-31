package domain

type Account struct {
	Number   string
	Document string
	Name     string
	Balance  int64
}

func NewAccount(number, document, name string) *Account {
	return &Account{
		Number:   number,
		Document: document,
		Name:     name,
	}
}
