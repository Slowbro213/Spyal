package events

import (
	"spyal/contracts"
)

type LeftEvent struct {
	BaseEvent
}

func NewLeftEvent(data map[string]any) contracts.Event {
	eventData := map[string]any{
		"type": Leftevent,
		"msg":  data,
	}
	return &LeftEvent{
		BaseEvent: BaseEvent{
			Data: eventData,
		},
	}
}

func (ee *LeftEvent) GetName() contracts.EventName {
	return Leftevent
}

func (ee *LeftEvent) Channel() string {
	return "game"
}

func (ee *LeftEvent) Topic() string {
	v, _ := ee.Data["topic"].(string)
	return v
}
