package worker

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/kpacha/marathon-pipeline/pipeline"
)

func ExampleWebhookConsume() {
	ts := creteTestServer("ok")
	defer ts.Close()

	eventType := "deployment_.*"
	appId := "group/.*"
	task := pipeline.Task{
		ID:   "id",
		Name: "name",
		Params: map[string]string{
			"worker":  "webhook",
			"method":  "POST",
			"url":     fmt.Sprintf("%s/some-uri", ts.URL),
			"payload": `payload={"text": "{{.ID}}: {{.Type}} [{{.Status}} at {{.Node}}]."}`,
		},
		Filter: &pipeline.FilterConstraint{
			EventType: &eventType,
			AppId:     &appId,
		},
	}
	webhookFactory := WebhookFactory{}
	webhook, err := webhookFactory.Build(task)
	fmt.Println(err)

	job := &pipeline.MarathonEvent{
		Type:   "test1",
		Status: "status1",
		ID:     "group/app1",
		Node:   "someNode",
	}
	fmt.Println(webhook.Consume(job))

	// Output:
	// <nil>
	// POST
	// /some-uri
	// payload={"text": "group/app1: test1 [status1 at someNode]."}
	// <nil>
}

func creteTestServer(msg string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method)
		fmt.Println(r.RequestURI)
		body, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(body))
		fmt.Fprintln(w, msg)
	}))
}
