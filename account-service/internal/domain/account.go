package domain

import (
	"errors"
	"fmt"
	"time"
	"unicode/utf8"
)

const (
	MinimumLengthName = 5
	MaximumLengthName = 120

	CPFLength  = 11
	CNPJLength = 14
)

type Account struct {
	Id        string
	Number    string
	Name      string
	Document  string
	Balance   int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewAccount(number string, document string, name string) *Account {
	return &Account{
		Number:    number,
		Document:  document,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Balance:   0,
	}
}

func (acc *Account) Validate() error {
	nameLength := utf8.RuneCountInString(acc.Name)
	if nameLength < MinimumLengthName || nameLength > MaximumLengthName {
		return fmt.Errorf("invalid name, should be between %v and %v characters", MinimumLengthName, MaximumLengthName)
	}

	documentLength := utf8.RuneCountInString(acc.Document)
	if documentLength != CPFLength && documentLength != CNPJLength {
		return fmt.Errorf("invalid document, should be CPF with %v or CNPJ with %v characters", CPFLength, CNPJLength)
	}

	return nil
}

func (acc *Account) Deposit(value int64) error {
	if value < 0 {
		return errors.New("for a deposit the value must be greater than zero")
	}

	acc.Balance += value
	acc.UpdatedAt = time.Now()

	return nil
}
