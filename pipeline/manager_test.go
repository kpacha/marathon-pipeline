package pipeline

import (
	"fmt"
	"time"
)

func ExampleManager() {
	wf := newTestWorkerFactory()
	store := NewMemoryTaskStore()
	events := make(chan *MarathonEvent)
	manager := NewManager(wf, &store, events)

	var errors []error
	go func() {
		for err := range manager.Error {
			errors = append(errors, err)
		}
	}()

	task1 := Task{
		ID:     "id1",
		Name:   "name1",
		Params: map[string]string{},
	}

	testWorkerParams := map[string]string{
		"worker":  "test1",
		"payload": "random text",
	}

	task2 := Task{
		ID:     "id2",
		Name:   "name2",
		Params: testWorkerParams,
	}

	eventType := "deployment_.*"
	appId := "group/.*"
	filter := FilterConstraint{
		EventType: &eventType,
		AppId:     &appId,
	}

	task3 := Task{
		ID:     "id3",
		Name:   "name3",
		Params: testWorkerParams,
		Filter: &filter,
	}

	time.Sleep(5 * time.Millisecond)

	store.Set(task1)
	store.Set(task2)
	store.Set(task3)

	time.Sleep(5 * time.Millisecond)

	events <- &MarathonEvent{Type: "deployment_failed", ID: "group/appA"}
	events <- &MarathonEvent{Type: "deployment_success", ID: "appB"}
	events <- &MarathonEvent{Type: "app_terminated_event", Status: "TASK_FINISHED", ID: "group/appC"}
	events <- &MarathonEvent{Type: "deployment_success", ID: "group/appD"}

	time.Sleep(5 * time.Millisecond)

	for _, err := range errors {
		fmt.Println(err)
	}

	// Output:
	// WARNING: malformed task! ProxyWorkerFactory: undefined param worker {id1 name1 <nil> map[]}
	// WARNING: malformed task! ProxyWorkerFactory: undefined param worker {id1 name1 <nil> map[]}
	// WARNING: malformed task! ProxyWorkerFactory: undefined param worker {id1 name1 <nil> map[]}
	// &{deployment_failed  group/appA 0001-01-01 00:00:00 +0000 UTC }
	// &{deployment_failed  group/appA 0001-01-01 00:00:00 +0000 UTC }
	// &{deployment_success  appB 0001-01-01 00:00:00 +0000 UTC }
	// &{app_terminated_event TASK_FINISHED group/appC 0001-01-01 00:00:00 +0000 UTC }
	// &{deployment_success  group/appD 0001-01-01 00:00:00 +0000 UTC }
	// &{deployment_success  group/appD 0001-01-01 00:00:00 +0000 UTC }
}

func newTestWorkerFactory() WorkerFactory {
	return ProxyWorkerFactory{map[string]WorkerFactory{
		"test1": TestWorkerFactory{},
		"test2": TestWorkerFactory2{},
	}}
}

type TestWorkerFactory struct{}

func (t TestWorkerFactory) Build(task Task) (Worker, error) {
	return TestWorker{}, nil
}

type TestWorkerFactory2 struct{}

func (t TestWorkerFactory2) Build(task Task) (Worker, error) {
	return TestWorker2{}, nil
}
