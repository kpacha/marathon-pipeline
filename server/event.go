package server

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"time"
)

type MarathonEvent struct {
	Type   string `json:"eventType"`
	Status string `json:"taskStatus,omitempty"`
	ID     string `json:"appId,omitempty"`
	Time   time.Time
	Node   string
}

type MarathonEventsParser struct{}

func (mp MarathonEventsParser) Parse(r io.ReadCloser, from string) (*MarathonEvent, error) {
	defer r.Close()
	var event *MarathonEvent
	body, err := ioutil.ReadAll(r)
	if err != nil {
		log.Println("Error reading from", r)
		return nil, err
	}
	if err = json.Unmarshal(body, &event); err != nil {
		log.Println("Error parsing to MarathonEvents")
		return nil, err
	}
	event.Node = from
	event.Time = time.Now()
	return event, nil
}
