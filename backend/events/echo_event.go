package events

import (
	"spyal/contracts"
)


type EchoEvent struct {
	BaseEvent
}

func NewEchoEvent(data map[string]any) contracts.Event{
	return &EchoEvent{
		BaseEvent: BaseEvent{
			Data: data,
		},
	}
}


func (ee *EchoEvent) GetName() contracts.EventName {
	return Echoevent
}

func (ee *EchoEvent) Channel() string{
	return "echo"
}
