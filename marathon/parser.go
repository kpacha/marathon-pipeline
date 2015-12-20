package marathon

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"time"

	"github.com/kpacha/marathon-pipeline/pipeline"
)

type MarathonEventsParser struct{}

func (mp MarathonEventsParser) Parse(r io.ReadCloser, from string) (*pipeline.MarathonEvent, error) {
	defer r.Close()
	var event *pipeline.MarathonEvent
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
