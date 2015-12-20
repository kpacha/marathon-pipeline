package marathon

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type MarathonConfig struct {
	Marathon     []MarathonServer
	Host         string
	Port         int
	BufferSize   int
	CallbackPath string
	callback     string
}

type MarathonServer struct {
	Host         string
	Port         int
	subsEndpoint string
}

type MarathonSubscriber struct {
	config *MarathonConfig
	parser MarathonEventsParser
	Buffer chan *MarathonEvent
}

const (
	defaultBufferSize   = 1000
	defaultCallbackPath = "marathon-listener"
)

func NewMarathonSubscription(config *MarathonConfig, p MarathonEventsParser) chan *MarathonEvent {
	mes := NewMarathonSubscriber(config, p)
	go mes.run()
	return mes.Buffer
}

func NewMarathonSubscriber(config *MarathonConfig, p MarathonEventsParser) MarathonSubscriber {
	return MarathonSubscriber{
		config: sanitizeConfig(config),
		parser: p,
		Buffer: make(chan *MarathonEvent, config.BufferSize),
	}
}

func sanitizeConfig(config *MarathonConfig) *MarathonConfig {
	if config.BufferSize == 0 {
		config.BufferSize = defaultBufferSize
	}
	if config.CallbackPath == "" {
		config.CallbackPath = defaultCallbackPath
	}
	config.CallbackPath = fmt.Sprintf("/%s", strings.Trim(config.CallbackPath, "/"))
	config.Host = strings.Trim(config.Host, "/")
	config.callback = fmt.Sprintf("http://%s:%d%s", config.Host, config.Port, config.CallbackPath)
	for s, server := range config.Marathon {
		config.Marathon[s].subsEndpoint = fmt.Sprintf("http://%s:%d/v2/eventSubscriptions", server.Host, server.Port)
	}
	return config
}

func (mes MarathonSubscriber) run() {
	mes.Register()
	defer mes.Unregister()

	http.HandleFunc(mes.config.CallbackPath, mes.MarathonListener)
	log.Println(http.ListenAndServe(fmt.Sprintf(":%d", mes.config.Port), nil))
	log.Println("DONE!")
}

func (mes MarathonSubscriber) Register() error {
	for _, server := range mes.config.Marathon {
		if !mes.isAlreadyRegistered(server.subsEndpoint) {
			if err := mes.registerEndpoint(server.subsEndpoint); err != nil {
				return err
			}
		}
	}
	return nil
}

func (mes MarathonSubscriber) Unregister() error {
	for _, server := range mes.config.Marathon {
		if err := mes.unregisterEndpoint(server.subsEndpoint); err != nil {
			return err
		}
	}
	return nil
}

func (mes MarathonSubscriber) MarathonListener(w http.ResponseWriter, r *http.Request) {
	defer func() { mes.Unregister() }()
	event, err := mes.parser.Parse(r.Body, r.RemoteAddr)
	if err == nil {
		mes.Buffer <- event
	}
}

func (mes MarathonSubscriber) registerEndpoint(subsEndpoint string) error {
	err := mes.sendHttpRequest("POST", subsEndpoint)
	if err != nil {
		return fmt.Errorf("Error registering: %s", err)
	}
	return nil

}

func (mes MarathonSubscriber) unregisterEndpoint(subsEndpoint string) error {
	err := mes.sendHttpRequest("DELETE", subsEndpoint)
	if err != nil {
		return fmt.Errorf("Error unregistering: %s", err)
	}
	return nil
}

func (mes MarathonSubscriber) sendHttpRequest(method, subsEndpoint string) error {
	endpoint := fmt.Sprintf("%s?callbackUrl=%s", subsEndpoint, mes.config.callback)

	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Error sending %s to %s. Status code [%d]", method, subsEndpoint, resp.StatusCode)
	}

	return nil
}

func (mes MarathonSubscriber) isAlreadyRegistered(subsEndpoint string) bool {
	resp, err := http.Get(subsEndpoint)
	if err != nil {
		log.Fatalf("Error retrieving event subscribers list: %s", err)
	}
	defer resp.Body.Close()

	callbacks := MarathonRedisteredCallbacks{}
	if err := json.NewDecoder(resp.Body).Decode(&callbacks); nil != err {
		log.Fatalf("Error decoding event subscribers list: %s", err)
	}

	for _, url := range callbacks.URL {
		if url == mes.config.callback {
			return true
		}
	}

	return false
}

type MarathonRedisteredCallbacks struct {
	URL []string `json:"callbackUrls"`
}
