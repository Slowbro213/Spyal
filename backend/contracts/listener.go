package contracts 



type Listener interface {
	GetEventName() EventName
	Handle(map[string]any)
}

type NewListenerFunc func() Listener
