package events

import "spyal/contracts"

type BaseEvent struct {
	Name   contracts.EventName
	Data   map[string]any
}

func NewBaseEvent(data map[string]any) contracts.Event {
	return &BaseEvent{
		Data: data,
	}
}

func (e *BaseEvent) GetName() contracts.EventName {
	return Baseevent
}

func (e *BaseEvent) GetData() map[string]any {
	return e.Data
}
