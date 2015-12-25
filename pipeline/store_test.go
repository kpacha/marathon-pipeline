package pipeline

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemoryTaskStore_GetAll(t *testing.T) {
	initialSnapshot := map[string]Task{"supu": Task{Name: "Name to display", ID: "supu"}}
	store, _, _ := newMemoryTaskStore(initialSnapshot)

	snapshot, err := store.GetAll()
	assert.Equal(t, initialSnapshot, snapshot)
	assert.Equal(t, nil, err)
}

func TestMemoryTaskStore_Get(t *testing.T) {
	taskId := "supu"
	initialSnapshot := map[string]Task{taskId: Task{Name: "Name to display", ID: taskId}}
	store, _, _ := newMemoryTaskStore(initialSnapshot)

	task, ok, err := store.Get(taskId)
	assert.Equal(t, initialSnapshot[taskId], task)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	task, ok, err = store.Get("unknown")
	assert.Equal(t, Task{}, task)
	assert.Equal(t, false, ok)
	assert.Equal(t, nil, err)
}

func TestMemoryTaskStore_Set(t *testing.T) {
	initialSnapshot := map[string]Task{}
	store, _, err := newMemoryTaskStore(initialSnapshot)
	store, subscription, err := newMemoryTaskStore(initialSnapshot)
	assert.Equal(t, nil, err)

	testTask := Task{Name: "Name to display", ID: "supu"}
	expectedSnapshot := map[string]Task{testTask.ID: testTask}

	task, ok, err := store.Get(testTask.ID)
	assert.Equal(t, Task{}, task)
	assert.Equal(t, false, ok)
	assert.Equal(t, nil, err)

	err = store.Set(testTask)
	assert.Equal(t, nil, err)

	select {
	case snapshot := <-subscription.SnapshotStream:
		assert.Equal(t, expectedSnapshot, snapshot)
	case <-time.After(5 * time.Millisecond):
		t.Error("timeout waiting for the snapshot after a set")
	}

	task, ok, err = store.Get(testTask.ID)
	assert.Equal(t, testTask, task)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)
}

func TestMemoryTaskStore_Delete(t *testing.T) {
	taskId := "supu"
	initialSnapshot := map[string]Task{taskId: Task{Name: "Name to display", ID: taskId}}
	store, subscription, err := newMemoryTaskStore(initialSnapshot)
	assert.Equal(t, nil, err)

	task, ok, err := store.Get(taskId)
	assert.Equal(t, initialSnapshot[taskId], task)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	err = store.Delete(taskId)
	assert.Equal(t, nil, err)

	select {
	case snapshot := <-subscription.SnapshotStream:
		assert.Equal(t, initialSnapshot, snapshot)
		assert.Equal(t, map[string]Task{}, snapshot)
	case <-time.After(5 * time.Millisecond):
		t.Error("timeout waiting for the snapshot after a delete")
	}

	task, ok, err = store.Get(taskId)
	assert.Equal(t, Task{}, task)
	assert.Equal(t, false, ok)
	assert.Equal(t, nil, err)
}

func newMemoryTaskStore(snapshot map[string]Task) (MemoryTaskStore, TaskStoreSubscription, error) {
	store := NewMemoryTaskStore()
	store.Store = snapshot
	subscription, err := store.Subscribe()
	return store, subscription, err
}
