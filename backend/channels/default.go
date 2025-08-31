package channels

import (
	"spyal/core"
	"spyal/contracts"
)

type Default struct{
	contracts.Channel
}

func NewDefaultChannel() contracts.Channel {
	return &Echo{
		Channel: core.NewChannel("default"),
	}
}
