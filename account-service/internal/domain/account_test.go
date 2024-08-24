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
