package events

import (
	"encoding/json"
	"errors"
	"reflect"
)

type EventPublish struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func NewEventPublish(event any) (*EventPublish, error) {
	if event == nil {
		return nil, errors.New("nil event input")
	}

	eventType := reflect.TypeOf(event)

	eventDataSerialized, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	return &EventPublish{
		Type: eventType.Elem().Name(),
		Data: string(eventDataSerialized),
	}, nil
}
