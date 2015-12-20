package pipeline

import (
	"fmt"
)

type Worker interface {
	Consume(job *MarathonEvent) error
}

type EventManager struct {
	Job     chan *MarathonEvent
	Workers []Worker
	Filters []FilterConstraint
	filters []*Filter
	Error   chan error
}

func NewEventManager(input chan *MarathonEvent, workers []Worker, filters []FilterConstraint) EventManager {
	em := EventManager{
		Job:     input,
		Workers: workers,
		Filters: filters,
		Error:   make(chan error, 1000),
	}
	em.Build()
	go em.Run()
	return em
}

func (em *EventManager) Build() {
	for _, f := range em.Filters {
		em.filters = append(em.filters, NewFilter(f))
	}
}

func (em *EventManager) Run() {
	for {
		job := <-em.Job
		if err := em.Consume(job); err != nil {
			em.Error <- err
		}
	}
}

func (em EventManager) Consume(job *MarathonEvent) error {
	var errors []error
	for step := range em.Workers {
		if err := em.consume(job, step); err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) != 0 {
		return errors[0]
	}
	return nil
}

func (em *EventManager) consume(job *MarathonEvent, w int) error {
	if w > len(em.Workers) {
		return fmt.Errorf("Worker id out of bounds")
	}
	if em.shouldConsume(job, w) {
		return em.Workers[w].Consume(job)
	}
	return nil
}

func (em *EventManager) shouldConsume(job *MarathonEvent, w int) bool {
	if w > len(em.filters) {
		return false
	}
	return em.filters[w].ShouldConsume(job)
}
