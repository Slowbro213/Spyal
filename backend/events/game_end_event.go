package events

import (
	"spyal/contracts"
)

type GameEndEvent struct {
	BaseEvent
}

func NewGameEndEvent(data map[string]any) contracts.Event {
	eventData := map[string]any{
		"type": Gameendevent,
		"msg":  data,
	}

	return &GameEndEvent{
		BaseEvent: BaseEvent{
			Data: eventData,
		},
	}
}

func (ee *GameEndEvent) GetName() contracts.EventName {
	return Gameendevent
}

func (ee *GameEndEvent) Channel() string {
	return "game"
}

func (ee *GameEndEvent) Topic() string {
	v, _ := ee.Data["topic"].(string)
	return v
}
