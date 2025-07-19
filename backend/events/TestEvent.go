package events

import "spyal/contracts"

type TestEvent struct {
	Event
}

func NewTestEvent(data map[string]any) contracts.Eventer{
	return &EchoEvent{
		Event: Event{
			Name: Testevent,
			Data: data,
		},
	}
}

