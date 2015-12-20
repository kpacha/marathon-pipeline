package worker

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/kpacha/marathon-pipeline/pipeline"
)

type Webhook struct {
	URL     []string
	Method  string
	Payload string
}

func (w Webhook) Consume(job *pipeline.MarathonEvent) error {
	client := &http.Client{}
	payload, err := w.parsePayload(job)
	if err != nil {
		return err
	}
	for _, url := range w.URL {
		req, err := http.NewRequest(w.Method, url, bytes.NewBuffer(payload))
		if err != nil {
			return err
		}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		resp.Body.Close()
	}
	return nil
}

func (w Webhook) parsePayload(job *pipeline.MarathonEvent) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})

	tmpl, err := template.New("test").Parse(w.Payload)
	if err != nil {
		return []byte{}, err
	}
	err = tmpl.Execute(buf, job)

	return buf.Bytes(), nil
}
