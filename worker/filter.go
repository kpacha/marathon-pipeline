package worker

import (
	"regexp"

	"github.com/kpacha/marathon-pipeline/marathon"
)

const wildCard = ".*"

type FilterConstraint struct {
	EventType  *string
	TaskStatus *string
	AppId      *string
}

type Filter struct {
	Constraint FilterConstraint
}

func NewFilter(f FilterConstraint) *Filter {
	return &Filter{f}
}

func (f *Filter) ShouldConsume(job *marathon.MarathonEvent) bool {
	return f.check(f.Constraint.EventType, job.Type) &&
		f.check(f.Constraint.TaskStatus, job.Status) &&
		f.check(f.Constraint.AppId, job.ID)
}

func (f *Filter) check(pattern *string, value string) bool {
	match, _ := regexp.MatchString(f.getSanitized(pattern), value)
	return match
}

func (f *Filter) getSanitized(s *string) string {
	if s != nil && *s != "" {
		return *s
	}
	return wildCard
}
