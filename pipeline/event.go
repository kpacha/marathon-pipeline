package pipeline

import "time"

type MarathonEvent struct {
	Type   string `json:"eventType"`
	Status string `json:"taskStatus,omitempty"`
	ID     string `json:"appId,omitempty"`
	Time   time.Time
	Node   string
}

type Task struct {
	ID     string
	Name   string
	Filter *FilterConstraint
	Params map[string]string
}

type FilterConstraint struct {
	EventType  *string
	TaskStatus *string
	AppId      *string
}
