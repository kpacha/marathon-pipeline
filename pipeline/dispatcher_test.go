package pipeline

import (
	"fmt"
	"time"
)

var (
	emTestJob1 = &MarathonEvent{
		Type:   "test1",
		Status: "status1",
		ID:     "group/app1",
	}
	emTestJob2 = &MarathonEvent{
		Type:   "test2",
		Status: "status2",
		ID:     "group/app2",
	}
	emTestJob3 = &MarathonEvent{
		Type:   "test3",
		Status: "status3",
		ID:     "app3",
	}
)

func ExampleFilterGeneration() {
	input := make(chan *MarathonEvent)

	NewDispatcher(input, []Worker{newTestWorker()})

	testDispatchingSequence(input)

	time.Sleep(time.Millisecond)

	// Output:
	// &{test1 status1 group/app1 0001-01-01 00:00:00 +0000 UTC }
	// &{test2 status2 group/app2 0001-01-01 00:00:00 +0000 UTC }
	// &{test1 status1 group/app1 0001-01-01 00:00:00 +0000 UTC }
}

func ExampleErrorHandling() {
	input := make(chan *MarathonEvent)

	dispatcher := NewDispatcher(input, []Worker{newTestWorker(), newTestFailWorker()})

	var errors []error
	go func() {
		for err := range dispatcher.Error {
			errors = append(errors, err)
		}
	}()

	testDispatchingSequence(input)

	time.Sleep(time.Millisecond)

	for _, err := range errors {
		fmt.Println(err)
	}

	// Output:
	// &{test1 status1 group/app1 0001-01-01 00:00:00 +0000 UTC }
	// &{test2 status2 group/app2 0001-01-01 00:00:00 +0000 UTC }
	// &{test1 status1 group/app1 0001-01-01 00:00:00 +0000 UTC }
	// TestError: &{test3 status3 app3 0001-01-01 00:00:00 +0000 UTC }
	// TestError: &{test3 status3 app3 0001-01-01 00:00:00 +0000 UTC }
	// TestError: &{test1 status1 group/app1 0001-01-01 00:00:00 +0000 UTC }
	// TestError: &{test3 status3 app3 0001-01-01 00:00:00 +0000 UTC }
	// TestError: &{test2 status2 group/app2 0001-01-01 00:00:00 +0000 UTC }
	// TestError: &{test3 status3 app3 0001-01-01 00:00:00 +0000 UTC }
	// TestError: &{test1 status1 group/app1 0001-01-01 00:00:00 +0000 UTC }
}

func ExampleWorkerReplacement() {
	input := make(chan *MarathonEvent)

	wrapper1 := newTestWorker()
	wrapper2 := newTestFailWorker()
	dispatcher := NewDispatcher(input, []Worker{wrapper1})

	var errors []error
	go func() {
		for err := range dispatcher.Error {
			errors = append(errors, err)
		}
	}()

	testDispatchingSequence(input)

	time.Sleep(5 * time.Millisecond)

	dispatcher.workerStream <- []Worker{wrapper1, wrapper2}

	time.Sleep(time.Millisecond)

	testDispatchingSequence(input)

	time.Sleep(5 * time.Millisecond)

	for _, err := range errors {
		fmt.Println(err)
	}

	// Output:
	// &{test1 status1 group/app1 0001-01-01 00:00:00 +0000 UTC }
	// &{test2 status2 group/app2 0001-01-01 00:00:00 +0000 UTC }
	// &{test1 status1 group/app1 0001-01-01 00:00:00 +0000 UTC }
	// &{test1 status1 group/app1 0001-01-01 00:00:00 +0000 UTC }
	// &{test2 status2 group/app2 0001-01-01 00:00:00 +0000 UTC }
	// &{test1 status1 group/app1 0001-01-01 00:00:00 +0000 UTC }
	// TestError: &{test3 status3 app3 0001-01-01 00:00:00 +0000 UTC }
	// TestError: &{test3 status3 app3 0001-01-01 00:00:00 +0000 UTC }
	// TestError: &{test1 status1 group/app1 0001-01-01 00:00:00 +0000 UTC }
	// TestError: &{test3 status3 app3 0001-01-01 00:00:00 +0000 UTC }
	// TestError: &{test2 status2 group/app2 0001-01-01 00:00:00 +0000 UTC }
	// TestError: &{test3 status3 app3 0001-01-01 00:00:00 +0000 UTC }
	// TestError: &{test1 status1 group/app1 0001-01-01 00:00:00 +0000 UTC }
}

func newTestWorker() Worker {
	taskPattern := "status\\d+"
	appPattern := "group/.*"
	fc := FilterConstraint{TaskStatus: &taskPattern, AppId: &appPattern}
	return NewWorkerWrapper(TestWorker{}, &fc)
}

func newTestFailWorker() Worker {
	return NewWorkerWrapper(TestWorker2{}, nil)
}

func testDispatchingSequence(input chan *MarathonEvent) {
	input <- emTestJob3
	input <- emTestJob3
	input <- emTestJob1
	input <- emTestJob3
	input <- emTestJob2
	input <- emTestJob3
	input <- emTestJob1
}

type TestWorker struct{}

func (t TestWorker) Consume(job *MarathonEvent) error {
	fmt.Println(job)
	return nil
}

type TestWorker2 struct{}

func (t TestWorker2) Consume(job *MarathonEvent) error {
	return fmt.Errorf("TestError: %v", job)
}
