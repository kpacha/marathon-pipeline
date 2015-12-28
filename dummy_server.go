package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/kpacha/marathon-pipeline/marathon"
	"github.com/kpacha/marathon-pipeline/pipeline"
	"github.com/kpacha/marathon-pipeline/server"
	"github.com/kpacha/marathon-pipeline/worker"
	"github.com/kpacha/marathon-pipeline/zookeeper"
)

func main() {
	zkHost := flag.String("zk", "locahost:2181", "hostname")
	host := flag.String("h", "locahost", "hostname")
	listenerPort := flag.Int("e", 8081, "listener port")
	restPort := flag.Int("p", 8080, "rest API port")
	ttl := flag.Duration("l", 60*time.Second, "time to live of the app")
	flag.Parse()

	subscriber := geteventSubscriber(*host, *listenerPort)
	defer subscriber.Unregister()

	zk := startZK(*zkHost)

	manager := pipeline.NewManager(initWorkerFactory(), &zk, subscriber.Buffer)

	startRestServer(*restPort, zk)

	go func() {
		for err := range manager.Error {
			fmt.Println("error:", err)
		}
	}()

	time.Sleep(*ttl)
}

func geteventSubscriber(host string, port int) marathon.MarathonSubscriber {
	config := &marathon.MarathonConfig{
		Marathon: []marathon.MarathonServer{marathon.MarathonServer{Host: "marathon.mesos", Port: 8080}},
		Host:     host,
		Port:     port,
	}
	return marathon.NewMarathonSubscriber(config, marathon.MarathonEventsParser{})
}

func startZK(zkHost string) pipeline.TaskStore {
	zk, err := zookeeper.NewZKMemoryTaskStore([]string{zkHost})
	if err != nil {
		panic(err.Error())
	}
	return zk
}

func startRestServer(port int, store pipeline.TaskStore) {
	restServer := server.Default(port, store)
	go restServer.Run()
}

func initWorkerFactory() pipeline.WorkerFactory {
	return pipeline.ProxyWorkerFactory{
		Factories: map[string]pipeline.WorkerFactory{
			"webhook": worker.WebhookFactory{},
		},
	}
}
