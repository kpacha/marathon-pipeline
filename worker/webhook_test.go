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

	webhook := Webhook{
		URL:     []string{fmt.Sprintf("%s/some-uri", ts.URL)},
		Method:  "GET",
		Payload: "{\"message\":\"{{.ID}}: {{.Status}} [{{.Type}}]\"}",
	}

	job := &pipeline.MarathonEvent{
		Type:   "test1",
		Status: "status1",
		ID:     "group/app1",
	}
	if err := webhook.Consume(job); err != nil {
		panic(err)
	}

	// Output:
	// GET
	// /some-uri
	// {"message":"group/app1: status1 [test1]"}
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
