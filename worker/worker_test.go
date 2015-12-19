package worker

import (
	"fmt"
	"time"

	"github.com/kpacha/marathon-pipeline/server"
)

var (
	emTestJob1 = &server.MarathonEvent{
		Type:   "test1",
		Status: "status1",
		ID:     "group/app1",
	}
	emTestJob2 = &server.MarathonEvent{
		Type:   "test2",
		Status: "status2",
		ID:     "group/app2",
	}
	emTestJob3 = &server.MarathonEvent{
		Type:   "test3",
		Status: "status3",
		ID:     "app3",
	}
)

func ExampleFilterGeneration() {
	taskPattern := "status\\d+"
	appPattern := "group/.*"
	fc := FilterConstraint{TaskStatus: &taskPattern, AppId: &appPattern}

	input := make(chan *server.MarathonEvent)

	NewEventManager(input, []Worker{TestWorker{}}, []FilterConstraint{fc})

	input <- emTestJob3
	input <- emTestJob3
	input <- emTestJob1
	input <- emTestJob3
	input <- emTestJob2
	input <- emTestJob3
	input <- emTestJob1

	time.Sleep(time.Millisecond)

	// Output:
	// &{test1 status1 group/app1 0001-01-01 00:00:00 +0000 UTC }
	// &{test2 status2 group/app2 0001-01-01 00:00:00 +0000 UTC }
	// &{test1 status1 group/app1 0001-01-01 00:00:00 +0000 UTC }
}

func ExampleErrorHandling() {
	taskPattern := "status\\d+"
	appPattern := "group/.*"
	fc := FilterConstraint{TaskStatus: &taskPattern, AppId: &appPattern}
	input := make(chan *server.MarathonEvent)

	em := NewEventManager(
		input,
		[]Worker{TestWorker{}, TestWorker2{}},
		[]FilterConstraint{fc, FilterConstraint{}})

	var errors []error
	go func() {
		for err := range em.Error {
			errors = append(errors, err)
		}
	}()

	input <- emTestJob1
	input <- emTestJob2
	input <- emTestJob3

	time.Sleep(time.Millisecond)

	for _, err := range errors {
		fmt.Println(err)
	}

	// Output:
	// &{test1 status1 group/app1 0001-01-01 00:00:00 +0000 UTC }
	// &{test2 status2 group/app2 0001-01-01 00:00:00 +0000 UTC }
	// TestError: &{test1 status1 group/app1 0001-01-01 00:00:00 +0000 UTC }
	// TestError: &{test2 status2 group/app2 0001-01-01 00:00:00 +0000 UTC }
	// TestError: &{test3 status3 app3 0001-01-01 00:00:00 +0000 UTC }
}

type TestWorker struct{}

func (t TestWorker) Consume(job *server.MarathonEvent) error {
	fmt.Println(job)
	return nil
}

type TestWorker2 struct{}

func (t TestWorker2) Consume(job *server.MarathonEvent) error {
	return fmt.Errorf("TestError: %v", job)
}
