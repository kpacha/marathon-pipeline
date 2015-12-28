package zookeeper

import (
	"fmt"
	"sort"
	"time"

	"github.com/kpacha/marathon-pipeline/pipeline"
)

func ExampleZKTaskStore() {
	testingRootPath := fmt.Sprintf("%s-testing", rootPath)
	rootPath = testingRootPath

	store, err := NewZKTaskStore([]string{"localhost:2181"})
	fmt.Println(err)

	task1 := pipeline.Task{
		ID:     "id1",
		Name:   "name1",
		Params: map[string]string{},
	}
	task2 := pipeline.Task{
		ID:     "id2",
		Name:   "name2",
		Params: map[string]string{},
	}
	task3 := pipeline.Task{
		ID:     "id3",
		Name:   "name3",
		Params: map[string]string{},
	}

	store.conn.Create(testingRootPath, []byte{}, int32(0), store.acl)

	subscription, err2 := store.Subscribe()
	fmt.Println(err2)

	dumpSorted(store.GetAll())
	fmt.Println(store.Set(task1))
	dumpSorted(store.GetAll())

	fmt.Println(store.Get(task1.ID))
	fmt.Println(store.Get(task2.ID))

	fmt.Println(store.Set(task2))
	dumpSorted(store.GetAll())

	time.Sleep(5 * time.Millisecond)

	fmt.Println(store.Set(task3))
	dumpSorted(store.GetAll())

	time.Sleep(5 * time.Millisecond)

	fmt.Println(store.Delete(task1.ID))
	dumpSorted(store.GetAll())

	time.Sleep(5 * time.Millisecond)

	fmt.Println(store.Delete(task2.ID))
	dumpSorted(store.GetAll())

	time.Sleep(5 * time.Millisecond)

	fmt.Println(store.Delete(task3.ID))
	dumpSorted(store.GetAll())

	for i := 0; i < 6; i++ {
		select {
		case snapshot := <-subscription.SnapshotStream:
			dumpSorted(snapshot, nil)
		case <-time.After(5 * time.Millisecond):
			fmt.Println("timeout waiting for the snapshot")
		}
	}

	store.conn.Delete(testingRootPath, 0)
	// Output:
	// <nil>
	// <nil>
	// []
	// <nil>
	// [{id1 name1 <nil> map[]}]
	// {id1 name1 <nil> map[]} true <nil>
	// {  <nil> map[]} false zk: node does not exist
	// <nil>
	// [{id1 name1 <nil> map[]} {id2 name2 <nil> map[]}]
	// <nil>
	// [{id1 name1 <nil> map[]} {id2 name2 <nil> map[]} {id3 name3 <nil> map[]}]
	// <nil>
	// [{id2 name2 <nil> map[]} {id3 name3 <nil> map[]}]
	// <nil>
	// [{id3 name3 <nil> map[]}]
	// <nil>
	// []
	// [{id1 name1 <nil> map[]}]
	// [{id1 name1 <nil> map[]} {id2 name2 <nil> map[]}]
	// [{id1 name1 <nil> map[]} {id2 name2 <nil> map[]} {id3 name3 <nil> map[]}]
	// [{id2 name2 <nil> map[]} {id3 name3 <nil> map[]}]
	// [{id3 name3 <nil> map[]}]
	// []
}

func dumpSorted(snapshot map[string]pipeline.Task, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
	var keys []string
	for _, k := range snapshot {
		keys = append(keys, k.ID)
	}
	sort.Strings(keys)
	var tasks []pipeline.Task
	for _, k := range keys {
		tasks = append(tasks, snapshot[k])
	}
	fmt.Println(tasks)
}
