package pipeline

type TaskStore interface {
	GetAll() (map[string]Task, error)
	Get(k string) (Task, bool, error)
	Set(t Task) error
	Delete(k string) error
	Subscribe() (TaskStoreSubscription, error)
}

type TaskStoreSubscription struct {
	SnapshotStream chan map[string]Task
	ErrorStream    chan error
}

type MemoryTaskStore struct {
	Store         map[string]Task
	Subscriptions *[]TaskStoreSubscription
}

func (m MemoryTaskStore) GetAll() (map[string]Task, error) {
	return m.Store, nil
}

func (m MemoryTaskStore) Get(k string) (Task, bool, error) {
	v, ok := m.Store[k]
	return v, ok, nil
}

func (m MemoryTaskStore) Set(t Task) error {
	m.Store[t.ID] = t
	m.notifySubscriptors()
	return nil
}

func (m MemoryTaskStore) Delete(k string) error {
	delete(m.Store, k)
	m.notifySubscriptors()
	return nil
}

func (m MemoryTaskStore) Subscribe() (TaskStoreSubscription, error) {
	susbscription := TaskStoreSubscription{
		make(chan map[string]Task, 1000),
		make(chan error, 1000),
	}
	*m.Subscriptions = append(*m.Subscriptions, susbscription)
	return susbscription, nil
}

func (m MemoryTaskStore) notifySubscriptors() {
	if m.Subscriptions == nil {
		return
	}
	for _, subscriptor := range *m.Subscriptions {
		go m.notify(subscriptor.SnapshotStream, m.Store)
	}
}

func (m MemoryTaskStore) notify(channel chan map[string]Task, snapshot map[string]Task) {
	channel <- snapshot
}
