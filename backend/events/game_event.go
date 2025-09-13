package events

import (
	"spyal/contracts"
)


type GameEvent struct {
	BaseEvent
}

func NewGameEvent(data map[string]any) contracts.Event{
	return &GameEvent{
		BaseEvent: BaseEvent{
			Data: data,
		},
	}
}

func (ee *GameEvent) GetName() contracts.EventName {
	return Gameevent
}

func (ee *GameEvent) Channel() string{
	return "game"
}

func (ee *GameEvent) Topic() string {
	v, _ := ee.Data["topic"].(string)
	return v
}

