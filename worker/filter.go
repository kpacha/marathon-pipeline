package worker

import (
	"regexp"

	"github.com/kpacha/marathon-pipeline/pipeline"
)

const wildCard = ".*"

type Filter struct {
	Constraint pipeline.FilterConstraint
}

func NewFilter(f pipeline.FilterConstraint) *Filter {
	return &Filter{f}
}

func (f *Filter) ShouldConsume(job *pipeline.MarathonEvent) bool {
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
