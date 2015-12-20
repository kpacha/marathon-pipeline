package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/kpacha/marathon-pipeline/marathon"
	"github.com/kpacha/marathon-pipeline/pipeline"
	"github.com/kpacha/marathon-pipeline/worker"
)

func main() {
	slackUrl := flag.String("s", "", "slack url")
	host := flag.String("h", "locahost", "hostname")
	port := flag.Int("p", 8080, "port")
	ttl := flag.Duration("l", 60*time.Second, "time to live of the app")
	flag.Parse()

	config := &marathon.MarathonConfig{
		Marathon: []marathon.MarathonServer{marathon.MarathonServer{Host: "marathon.mesos", Port: 8080}},
		Host:     *host,
		Port:     *port,
	}
	subscriber := marathon.NewMarathonSubscriber(config, marathon.MarathonEventsParser{})

	taskPattern := "deployment_.*"
	appPattern := "group/.*"
	fc := pipeline.FilterConstraint{TaskStatus: &taskPattern, AppId: &appPattern}

	webhook := worker.Webhook{
		URL:    []string{*slackUrl},
		Method: "POST",
		Payload: `payload={
				"username": "new-bot-name",
				"icon_emoji": ":robot_face:".
				"text": "{{.ID}}: {{.Type}} [{{.Status}} at {{.Node}}].\n<http://marathon.mesos:8080/v2/apps/{{.ID}}|Click here> for details!"
			}`,
	}

	em := pipeline.NewEventManager(subscriber.Buffer, []pipeline.Worker{webhook}, []pipeline.FilterConstraint{fc})

	go func() {
		for err := range em.Error {
			fmt.Println("error:", err)
		}
	}()

	time.Sleep(*ttl)
	subscriber.Unregister()
}
