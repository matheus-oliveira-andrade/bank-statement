package events

type EventPublish struct {
	Type string `json:"type"`
	Data string `json:"data"`
}
