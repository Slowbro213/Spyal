package contracts 



type Listener interface {
	GetEventName() EventName
	Handle(Event)
}

type NewListenerFunc func() Listener
