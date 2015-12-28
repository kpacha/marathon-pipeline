package pipeline

import "fmt"

type Worker interface {
	Consume(job *MarathonEvent) error
}

type WorkerWrapper struct {
	worker Worker
	filter Filter
}

func NewWorkerWrapper(w Worker, fc *FilterConstraint) Worker {
	f := Filter{}
	if fc != nil {
		f = *NewFilter(*fc)
	}
	return WorkerWrapper{w, f}
}

func (w WorkerWrapper) Consume(job *MarathonEvent) error {
	if w.filter.ShouldConsume(job) {
		return w.worker.Consume(job)
	}
	return nil
}

type WorkerFactory interface {
	Build(task Task) (Worker, error)
}

type ProxyWorkerFactory struct {
	Factories map[string]WorkerFactory
}

func (pf ProxyWorkerFactory) Build(task Task) (Worker, error) {
	name, ok := task.Params["worker"]
	if !ok {
		return nil, fmt.Errorf("ProxyWorkerFactory: undefined param worker")
	}
	f, ok := pf.Factories[name]
	if !ok {
		return nil, fmt.Errorf("ProxyWorkerFactory: factory not defined for workers of the type [%s]", name)
	}
	return f.Build(task)
}
