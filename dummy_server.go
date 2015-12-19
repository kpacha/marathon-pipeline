package main

import (
	"fmt"

	"github.com/kpacha/marathon-pipeline/server"
	"github.com/kpacha/marathon-pipeline/worker"
)

func main() {
	config := &server.MarathonConfig{
		Marathon: []server.MarathonServer{server.MarathonServer{Host: "marathon.mesos", Port: 8080}},
		Host:     "localhost",
		Port:     8080,
	}
	subscriber := server.NewMarathonSubscriber(config, server.MarathonEventsParser{})

	taskPattern := "status\\d+"
	appPattern := "group/.*"
	fc := worker.FilterConstraint{TaskStatus: &taskPattern, AppId: &appPattern}

	em := worker.NewEventManager(subscriber.Buffer, []worker.Worker{DummyWorker{}}, []worker.FilterConstraint{fc})

	go func() {
		for err := range em.Error {
			fmt.Println("error:", err)
		}
	}()

}

type DummyWorker struct{}

func (t DummyWorker) Consume(job *server.MarathonEvent) error {
	fmt.Println(job)
	return nil
}
