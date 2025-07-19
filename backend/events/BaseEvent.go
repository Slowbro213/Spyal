package events

import "spyal/contracts"

type Event struct {
	Name   contracts.EventName
	Data   map[string]any
}

func NewEvent(name contracts.EventName,data map[string]any) *Event {
	return &Event{
		Name: name,
		Data: data,
	}
}

func (e *Event) GetName() contracts.EventName {
	return e.Name
}

func (e *Event) GetData() map[string]any {
	return e.Data
}
