package events

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEventPublish_success(t *testing.T) {
	// Arrange
	event := NewAccountCreated("1", "name 1", "01234567890")
	expectedType := "AccountCreated"
	expectedData := `{"number":"1","name":"name 1","document":"01234567890"}`

	// Act
	result, err := NewEventPublish(event)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, result.Type, expectedType)
	assert.Equal(t, result.Data, expectedData)
}

func TestNewEventPublish_nillInput_error(t *testing.T) {
	// Arrange & Act
	result, err := NewEventPublish(nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestNewEventPublish_notSerializableInput_error(t *testing.T) {
	// Arrange
	event := make(chan int)

	// Act
	result, err := NewEventPublish(event)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}
