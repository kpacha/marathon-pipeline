package zookeeper

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/kpacha/marathon-pipeline/pipeline"
	"github.com/samuel/go-zookeeper/zk"
)

var rootPath = "/marathon-pipeline"

type ZKTaskStore struct {
	conn     *zk.Conn
	acl      []zk.ACL
	rootPath string
	mailer   *pipeline.TaskStoreSubscriptionHandler
}

func NewZKTaskStore(zks []string) (ZKTaskStore, error) {
	conn, _, err := zk.Connect(zks, time.Second)
	if err != nil {
		return ZKTaskStore{}, err
	}
	store := ZKTaskStore{
		conn:     conn,
		acl:      zk.WorldACL(zk.PermAll),
		rootPath: rootPath,
		mailer:   pipeline.NewTaskStoreSubscriptionHandler(),
	}
	store.init()
	return store, nil
}

func (zkStore ZKTaskStore) GetAll() (map[string]pipeline.Task, error) {
	children, _, err := zkStore.conn.Children(zkStore.rootPath)
	snapshot := map[string]pipeline.Task{}
	if err != nil {
		return snapshot, err
	}
	for _, name := range children {
		task, ok, err := zkStore.Get(name)
		if err != nil {
			return snapshot, err
		}
		if ok {
			snapshot[name] = task
		}
	}
	return snapshot, nil
}

func (zkStore ZKTaskStore) Get(k string) (pipeline.Task, bool, error) {
	var task pipeline.Task
	data, _, err := zkStore.conn.Get(zkStore.path(k))
	if err != nil {
		return task, false, err
	}
	if err := json.Unmarshal(data, &task); err != nil {
		return task, false, err
	}
	return task, true, nil
}

func (zkStore ZKTaskStore) Set(t pipeline.Task) error {
	tBytes, err := json.Marshal(t)
	if err != nil {
		return err
	}
	_, err = zkStore.conn.Create(zkStore.path(t.ID), tBytes, int32(0), zkStore.acl)
	return err
}

func (zkStore ZKTaskStore) Delete(k string) error {
	return zkStore.conn.Delete(zkStore.path(k), 0)
}

func (zkStore ZKTaskStore) Subscribe() (pipeline.TaskStoreSubscription, error) {
	return zkStore.mailer.Subscribe()
}

func (zkStore ZKTaskStore) path(k string) string {
	return fmt.Sprintf("%s/%s", zkStore.rootPath, k)
}

func (zkStore ZKTaskStore) init() error {
	_, _, err := zkStore.conn.Get(zkStore.rootPath)
	if err != nil {
		if err == zk.ErrNoNode {
			if _, err = zkStore.conn.Create(zkStore.rootPath, []byte{}, int32(0), zkStore.acl); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	go zkStore.advertise()
	return nil
}

func (zkStore ZKTaskStore) advertise() {
	var snapshot map[string]pipeline.Task
	for {
		if _, _, event, err := zkStore.conn.ChildrenW(zkStore.rootPath); err == nil {
			evt := <-event
			if evt.Type > 0 {
				if snapshot, err = zkStore.GetAll(); err == nil {
					zkStore.mailer.NotifySubscriptors(snapshot)
				}
			}
		}
	}
}
