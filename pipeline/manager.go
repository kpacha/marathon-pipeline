package pipeline

import "fmt"

type Manager struct {
	wf     WorkerFactory
	store  *TaskStore
	events chan *MarathonEvent
	Error  chan error
}

func NewManager(wf WorkerFactory, store *TaskStore, events chan *MarathonEvent) Manager {
	m := Manager{wf, store, events, make(chan error)}
	go m.Run()
	return m
}

func (m Manager) Run() {
	subscription, err := (*m.store).Subscribe()
	if err != nil {
		panic("fatal subscribing to the received subscribable!")
	}
	d := NewDispatcher(m.events, []Worker{})
	m.Error = d.Error
	for {
		snapshot := <-subscription.SnapshotStream
		d.workerStream <- m.getWorkers(snapshot)
	}
}

func (m Manager) getWorkers(snapshot map[string]Task) []Worker {
	workers := []Worker{}
	for _, task := range snapshot {
		w, err := m.wf.Build(task)
		if err == nil {
			workers = append(workers, NewWorkerWrapper(w, task.Filter))
		} else {
			fmt.Println("WARNING: malformed task!", err.Error(), task)
		}
	}
	return workers
}
