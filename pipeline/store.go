package pipeline

type TaskStore interface {
	GetAll() (map[string]Task, error)
	Get(k string) (Task, bool, error)
	Set(t Task) error
	Delete(k string) error
}

type MemoryTaskStore struct {
	Store map[string]Task
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
	return nil
}

func (m MemoryTaskStore) Delete(k string) error {
	delete(m.Store, k)
	return nil
}
