package contracts 

type Event interface {
	GetName() EventName
	GetData() map[string]any
}

type NewEventFunc func(data map[string]any) Event

type EventName int
