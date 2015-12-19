package main

import (
	"fmt"

	"github.com/kpacha/marathon-pipeline/server"
)

func main() {
	config := &server.MarathonConfig{
		Marathon: []server.MarathonServer{server.MarathonServer{Host: "marathon.mesos", Port: 8080}},
		Host:     "localhost",
		Port:     8080,
	}
	subscription := server.NewMarathonSubscription(config, server.MarathonEventsParser{})

	for {
		event := <-subscription
		fmt.Println(event)
	}
}
