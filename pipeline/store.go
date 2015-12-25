package pipeline

import "fmt"

type TaskStore interface {
	GetAll() (map[string]Task, error)
	Get(k string) (Task, bool, error)
	Set(t Task) error
	Delete(k string) error
	Subscribe() (TaskStoreSubscription, error)
}

type MemoryTaskStore struct {
	Store  map[string]Task
	mailer *TaskStoreSubscriptionHandler
}

func NewMemoryTaskStore() MemoryTaskStore {
	return MemoryTaskStore{
		Store:  map[string]Task{},
		mailer: NewTaskStoreSubscriptionHandler(),
	}
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
	if m.mailer == nil {
		return TaskStoreSubscription{}, fmt.Errorf("Undefined mailer for the memory task store")
	}
	return m.mailer.Subscribe()
}

func (m MemoryTaskStore) notifySubscriptors() {
	if m.mailer == nil {
		return
	}
	m.mailer.notifySubscriptors(m.Store)
}

type TaskStoreSubscription struct {
	SnapshotStream chan map[string]Task
	ErrorStream    chan error
}

type TaskStoreSubscriptionHandler struct {
	Subscriptions *[]TaskStoreSubscription
}

func NewTaskStoreSubscriptionHandler() *TaskStoreSubscriptionHandler {
	return &TaskStoreSubscriptionHandler{&[]TaskStoreSubscription{}}
}

func (sh TaskStoreSubscriptionHandler) Subscribe() (TaskStoreSubscription, error) {
	susbscription := TaskStoreSubscription{
		make(chan map[string]Task, 1000),
		make(chan error, 1000),
	}
	*sh.Subscriptions = append(*sh.Subscriptions, susbscription)
	return susbscription, nil
}

func (sh TaskStoreSubscriptionHandler) notifySubscriptors(snapshot map[string]Task) {
	if sh.Subscriptions == nil {
		return
	}
	for _, subscriptor := range *sh.Subscriptions {
		go sh.notify(subscriptor.SnapshotStream, snapshot)
	}
}

func (sh TaskStoreSubscriptionHandler) notify(channel chan map[string]Task, snapshot map[string]Task) {
	channel <- snapshot
}
