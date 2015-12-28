package pipeline

type Dispatcher struct {
	EventStream  chan *MarathonEvent
	workers      *[]Worker
	Error        chan error
	workerStream chan []Worker
}

func NewDispatcher(input chan *MarathonEvent, workers []Worker) *Dispatcher {
	d := Dispatcher{
		EventStream:  input,
		workers:      &workers,
		Error:        make(chan error, 1000),
		workerStream: make(chan []Worker),
	}
	go d.Run()
	return &d
}

func (d *Dispatcher) Run() {
	for {
		select {
		case job := <-d.EventStream:
			if err := d.consume(job); err != nil {
				d.Error <- err
			}
		case workers := <-d.workerStream:
			*d.workers = workers
		}
	}
}

func (d *Dispatcher) consume(job *MarathonEvent) error {
	var errors []error
	for _, w := range *d.workers {
		if err := w.Consume(job); err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) != 0 {
		return errors[0]
	}
	return nil
}

func (d Dispatcher) Consume(job *MarathonEvent) error {
	d.EventStream <- job
	return nil
}
