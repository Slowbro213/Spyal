package contracts

type ShouldBroadcast interface {
	Channel() string
	Topic()  string
}
