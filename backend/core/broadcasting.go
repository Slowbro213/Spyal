package core

import (
	"spyal/contracts"
)

type Channel struct {

}


type ShouldBroadcast interface {
	Channel(c *Channel)
}


//nolint
func Broadcast(e contracts.EventName, data map[string]any){

}
