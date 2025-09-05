package main

import (
	"log"
	"spyal/tooling/register"
)

const (
	listenersDir = "listeners"
	eventsDir = "events"
	channelsDir = "channels"
)

//gosec:disable G101 -- This is a false positive
func main() {
	dirs := [1]string{listenersDir}

	for _, dir := range dirs {
		constructors, err := register.GetListenerConstructors(dir)

		if err != nil {
			log.Fatalf("Error getting constructors %v", err)
		}

		err = register.GenerateListenerRegistry(dir, constructors)
		if err != nil {
			log.Fatalf("Error generating constructor registry %v", err)
		}
	}

	events, err := register.GetEventTypes(eventsDir)

	if err != nil {
		log.Fatalf("Error while getting events %v", err)
	}

	err = register.GenerateEventTypes(eventsDir, events)

	if err != nil {
		log.Fatalf("Error while generating events registry %v", err)
	}


	channels, err := register.GetChannelConstructors(channelsDir)

	if err != nil {
		log.Fatalf("Error while getting channels %v", err)
	}

	err = register.GenerateChannelRegistryFile(channelsDir,channels)

	if err != nil {
		log.Fatalf("Error while generating channels registry %v", err)
	}
}
