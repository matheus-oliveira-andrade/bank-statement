package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	// arrange
	number := "123456"
	document := "01234567890"
	name := "John Yos Bidden"

	// act
	acc := NewAccount(number, document, name)

	// assert
	assert.Equal(t, number, acc.Number)
	assert.Equal(t, document, acc.Document)
	assert.Equal(t, name, acc.Name)
	assert.Equal(t, int64(0), acc.Balance)
}

func TestAccountValidate(t *testing.T) {
	testCases := []struct {
		testName      string
		document      string
		name          string
		expectedError error
	}{
		{
			testName:      "given invalid name should return error about",
			document:      "01234567890",
			name:          "me",
			expectedError: errors.New("invalid name, should be between 5 and 120 characters"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {

			acc := NewAccount("", tc.document, tc.name)

			err := acc.Validate()

			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestDeposit_NegativeValue(t *testing.T) {
	// arrange
	acc := NewAccount("1", "01234567890", "John")
	acc.Balance = 0

	// act
	err := acc.Deposit(-1)

	// assert
	assert.Equal(t, err, errors.New("for a deposit the value must be greater than zero"))
	assert.Equal(t, int64(0), acc.Balance)
}

func TestDeposit_ValidValue(t *testing.T) {
	// arrange
	acc := NewAccount("1", "01234567890", "John")
	acc.Balance = 50

	// act
	err := acc.Deposit(100)

	// assert
	assert.Nil(t, err)
	assert.Equal(t, int64(150), acc.Balance)
}

func TestWithdraw_NegativeValue(t *testing.T) {
	// arrange
	initialBalance := int64(50)

	acc := NewAccount("1", "01234567890", "John")
	acc.Balance = initialBalance

	// act
	err := acc.withdraw(-10)

	// assert
	assert.NotNil(t, err)
	assert.Equal(t, initialBalance, acc.Balance)
}

func TestWithdraw_InsufficientFunds(t *testing.T) {
	// arrange
	initialBalance := int64(50)

	acc := NewAccount("1", "01234567890", "John")
	acc.Balance = initialBalance

	// act
	err := acc.withdraw(100)

	// assert
	assert.NotNil(t, err)
	assert.Equal(t, ErrInsufficientFunds, err)
	assert.Equal(t, initialBalance, acc.Balance)
}

func TestWithdraw_Success(t *testing.T) {
	// arrange
	initialBalance := int64(50)

	acc := NewAccount("1", "01234567890", "John")
	acc.Balance = initialBalance

	// act
	err := acc.withdraw(10)

	// assert
	assert.Nil(t, err)
	assert.Equal(t, initialBalance-10, acc.Balance)
}

func TestTransfer_ErrorWithdrawValueFromAccount(t *testing.T) {
	// arrange
	from := NewAccount("1", "01234567890", "John")
	from.Balance = 0

	to := NewAccount("2", "12345678901", "Jenny")

	// act
	err := from.Transfer(25, to)

	// assert
	assert.NotNil(t, ErrInsufficientFunds, err)
}

func TestTransfer_ErrorValueZero(t *testing.T) {
	// arrange
	from := NewAccount("1", "01234567890", "John")
	from.Balance = 0

	to := NewAccount("2", "12345678901", "Jenny")

	// act
	err := from.Transfer(0, to)

	// assert
	assert.NotNil(t, err)
}

func TestTransfer_Success(t *testing.T) {
	// arrange
	from := NewAccount("1", "01234567890", "John")
	from.Balance = 25

	to := NewAccount("2", "12345678901", "Jenny")
	to.Balance = 0

	// act
	err := from.Transfer(25, to)

	// assert
	assert.Nil(t, err)
	assert.Equal(t, int64(0), from.Balance)
	assert.Equal(t, int64(25), to.Balance)
}
