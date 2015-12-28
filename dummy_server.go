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

	config := &marathon.MarathonConfig{
		Marathon: []marathon.MarathonServer{marathon.MarathonServer{Host: "marathon.mesos", Port: 8080}},
		Host:     *host,
		Port:     *listenerPort,
	}
	subscriber := marathon.NewMarathonSubscriber(config, marathon.MarathonEventsParser{})

	workerFactory := pipeline.ProxyWorkerFactory{map[string]pipeline.WorkerFactory{
		"webhook": worker.WebhookFactory{},
	}}
	zk, err := zookeeper.NewZKMemoryTaskStore([]string{*zkHost})
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	manager := pipeline.NewManager(workerFactory, &zk, subscriber.Buffer)

	restServer := server.Default(*restPort, zk)
	go restServer.Run()

	go func() {
		for err := range manager.Error {
			fmt.Println("error:", err)
		}
	}()

	time.Sleep(*ttl)
	subscriber.Unregister()
}
