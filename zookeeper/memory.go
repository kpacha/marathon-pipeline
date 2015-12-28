package zookeeper

import "github.com/kpacha/marathon-pipeline/pipeline"

type ZKMemoryTaskStore struct {
	mem    *pipeline.TaskStore
	zk     *ZKTaskStore
	remote pipeline.TaskStoreSubscription
}

func NewZKMemoryTaskStore(zks []string) (ZKMemoryTaskStore, error) {
	mem := pipeline.NewMemoryTaskStore()
	zk, err := NewZKTaskStore(zks)
	if err != nil {
		return ZKMemoryTaskStore{}, err
	}
	zkSubscription, err := zk.Subscribe()
	if err != nil {
		return ZKMemoryTaskStore{}, err
	}
	zkMemoryTaskStore := ZKMemoryTaskStore{&mem, &zk, zkSubscription}
	zkMemoryTaskStore.init()
	return zkMemoryTaskStore, nil
}

func (zkm ZKMemoryTaskStore) GetAll() (map[string]pipeline.Task, error) {
	return (*zkm.mem).GetAll()
}

func (zkm ZKMemoryTaskStore) Get(k string) (pipeline.Task, bool, error) {
	return (*zkm.mem).Get(k)
}

func (zkm ZKMemoryTaskStore) Set(t pipeline.Task) error {
	return zkm.zk.Set(t)
}

func (zkm ZKMemoryTaskStore) Delete(k string) error {
	return zkm.zk.Delete(k)
}

func (zkm ZKMemoryTaskStore) Subscribe() (pipeline.TaskStoreSubscription, error) {
	return (*zkm.mem).Subscribe()
}

func (zkm ZKMemoryTaskStore) init() error {
	if snapshot, err := zkm.zk.GetAll(); err == nil {
		(*zkm.mem).Overwrite(snapshot)
	}
	go zkm.handleUpdates()
	return nil
}

func (zkm ZKMemoryTaskStore) handleUpdates() {
	for {
		snapshot := <-zkm.remote.SnapshotStream
		(*zkm.mem).Overwrite(snapshot)
	}
}
