package events

import "spyal/contracts"


type EchoEvent struct {
	Event
}

func NewEchoEvent(data map[string]any) contracts.Eventer{
	return &EchoEvent{
		Event: Event{
			Name: Echoevent,
			Data: data,
		},
	}
}
