package zookeeper

import (
	"fmt"
	"time"

	"github.com/kpacha/marathon-pipeline/pipeline"
)

func ExampleZKMemoryTaskStore_GetAll() {
	testingRootPath := fmt.Sprintf("%s-testing", rootPath)
	rootPath = testingRootPath

	store, err := NewZKMemoryTaskStore([]string{"localhost:2181"})
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

	store.zk.conn.Create(testingRootPath, []byte{}, int32(0), store.zk.acl)

	subscription, err2 := store.Subscribe()
	fmt.Println(err2)

	dumpSorted(store.GetAll())
	fmt.Println(store.Set(task1))

	time.Sleep(5 * time.Millisecond)

	data, _, err3 := store.zk.conn.Get(store.zk.path(task1.ID))
	fmt.Println(err3)
	fmt.Println(string(data))
	dumpSorted(store.GetAll())
	fmt.Println(store.Get(task1.ID))
	fmt.Println(store.Get(task2.ID))
	fmt.Println(store.Set(task2))

	time.Sleep(5 * time.Millisecond)

	dumpSorted(store.GetAll())
	fmt.Println(store.Set(task3))

	time.Sleep(5 * time.Millisecond)

	dumpSorted(store.GetAll())
	fmt.Println(store.Delete(task1.ID))

	time.Sleep(5 * time.Millisecond)

	dumpSorted(store.GetAll())
	fmt.Println(store.Delete(task2.ID))

	time.Sleep(5 * time.Millisecond)

	dumpSorted(store.GetAll())
	fmt.Println(store.Delete(task3.ID))

	time.Sleep(5 * time.Millisecond)

	dumpSorted(store.GetAll())

	for i := 0; i < 6; i++ {
		select {
		case snapshot := <-subscription.SnapshotStream:
			dumpSorted(snapshot, nil)
		case <-time.After(5 * time.Millisecond):
			fmt.Println("timeout waiting for the snapshot")
		}
	}

	store.zk.conn.Delete(testingRootPath, 0)
	// Output:
	// <nil>
	// <nil>
	// []
	// <nil>
	// <nil>
	// {"ID":"id1","Name":"name1","Filter":null,"Params":{}}
	// [{id1 name1 <nil> map[]}]
	// {id1 name1 <nil> map[]} true <nil>
	// {  <nil> map[]} false <nil>
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
