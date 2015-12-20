package marathon

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
)

func ExampleSubscriberSetsDefaultValuesForListening() {
	subscriber := createDefaultSubscriber()

	fmt.Println(subscriber.config.BufferSize)
	fmt.Println(subscriber.config.CallbackPath)
	fmt.Println(subscriber.config.callback)
	// Output:
	// 1000
	// /marathon-listener
	// http://example.com:9876/marathon-listener
}

func ExampleHandleEvent() {
	subscriber := createDefaultSubscriber()

	var jsonStr = []byte(`{"eventType":"typeA","taskStatus":"statusA","appId":"appA"}`)
	req, err := http.NewRequest("POST", "http://example.com:9876/marathon-listener", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	subscriber.MarathonListener(w, req)

	event := <-subscriber.Buffer

	fmt.Println(w.Code)
	fmt.Println(w.Body.String())
	fmt.Println(event.ID)
	fmt.Println(event.Node)
	fmt.Println(event.Status)
	fmt.Println(event.Type)
	// Output:
	// 200
	//
	// appA
	//
	// statusA
	// typeA
}

func ExampleRegister() {
	ts := creteTestServer("{\"callbackUrls\":[]}")
	defer ts.Close()

	subscriber := createSubscriberFromTestServer(ts.URL)

	if err := subscriber.Register(); err != nil {
		panic(err)
	}

	// Output:
	// GET
	// /v2/eventSubscriptions
	// []
	// POST
	// /v2/eventSubscriptions?callbackUrl=http://example.com:9876/marathon-listener
	// []
}

func ExampleAlreadyRegistered() {
	ts := creteTestServer("{\"callbackUrls\":[\"http://example.com:9876/marathon-listener\"]}")
	defer ts.Close()

	subscriber := createSubscriberFromTestServer(ts.URL)

	if err := subscriber.Register(); err != nil {
		panic(err)
	}

	// Output:
	// GET
	// /v2/eventSubscriptions
	// []
}

func ExampleUnregister() {
	ts := creteTestServer("")
	defer ts.Close()

	subscriber := createSubscriberFromTestServer(ts.URL)

	if err := subscriber.Unregister(); err != nil {
		panic(err)
	}

	// Output:
	// DELETE
	// /v2/eventSubscriptions?callbackUrl=http://example.com:9876/marathon-listener
	// []
}

func createDefaultSubscriber() MarathonSubscriber {
	return createSubscriber(MarathonServer{Host: "marathon.mesos", Port: 8080})
}

func createSubscriberFromTestServer(tsURL string) MarathonSubscriber {
	port, err := strconv.Atoi(strings.TrimPrefix(tsURL, "http://127.0.0.1:"))
	if err != nil {
		panic(err)
	}
	return createSubscriber(MarathonServer{Host: "127.0.0.1", Port: port})
}

func createSubscriber(s MarathonServer) MarathonSubscriber {
	config := &MarathonConfig{
		Marathon: []MarathonServer{s},
		Host:     "example.com",
		Port:     9876,
	}
	return NewMarathonSubscriber(config, MarathonEventsParser{})
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
		fmt.Println(body)
		fmt.Fprintln(w, msg)
	}))
}
