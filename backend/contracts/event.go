package contracts 

type Eventer interface {
	GetName() EventName
	GetData() map[string]any
}

type NewEventFunc func(data map[string]any) Eventer

type EventName int
