package worker

import (
	"fmt"

	"github.com/kpacha/marathon-pipeline/marathon"
)

var (
	filterTestJob1 = &marathon.MarathonEvent{
		Type:   "test1",
		Status: "status1",
		ID:     "group/app1",
	}
	filterTestJob2 = &marathon.MarathonEvent{
		Type:   "test2",
		Status: "status2",
		ID:     "group/app2",
	}
	filterTestJob3 = &marathon.MarathonEvent{
		Type:   "test3",
		Status: "status3",
		ID:     "app3",
	}
)

func ExampleEmptyFilter() {
	fc := FilterConstraint{}
	f := NewFilter(fc)
	fmt.Println(filterTestJob1)
	fmt.Println(f.ShouldConsume(filterTestJob1))
	fmt.Println(filterTestJob2)
	fmt.Println(f.ShouldConsume(filterTestJob2))
	fmt.Println(filterTestJob3)
	fmt.Println(f.ShouldConsume(filterTestJob3))
	// Output:
	// &{test1 status1 group/app1 0001-01-01 00:00:00 +0000 UTC }
	// true
	// &{test2 status2 group/app2 0001-01-01 00:00:00 +0000 UTC }
	// true
	// &{test3 status3 app3 0001-01-01 00:00:00 +0000 UTC }
	// true
}

func ExampleFilterByType() {
	typePattern := "test1"
	fc := FilterConstraint{EventType: &typePattern}
	f := NewFilter(fc)
	fmt.Println(filterTestJob1)
	fmt.Println(f.ShouldConsume(filterTestJob1))
	fmt.Println(filterTestJob2)
	fmt.Println(f.ShouldConsume(filterTestJob2))
	fmt.Println(filterTestJob3)
	fmt.Println(f.ShouldConsume(filterTestJob3))
	// Output:
	// &{test1 status1 group/app1 0001-01-01 00:00:00 +0000 UTC }
	// true
	// &{test2 status2 group/app2 0001-01-01 00:00:00 +0000 UTC }
	// false
	// &{test3 status3 app3 0001-01-01 00:00:00 +0000 UTC }
	// false
}

func ExampleComplexFilter() {
	taskPattern := "status\\d+"
	appPattern := "group/.*"
	fc := FilterConstraint{TaskStatus: &taskPattern, AppId: &appPattern}
	f := NewFilter(fc)
	fmt.Println(filterTestJob1)
	fmt.Println(f.ShouldConsume(filterTestJob1))
	fmt.Println(filterTestJob2)
	fmt.Println(f.ShouldConsume(filterTestJob2))
	fmt.Println(filterTestJob3)
	fmt.Println(f.ShouldConsume(filterTestJob3))
	// Output:
	// &{test1 status1 group/app1 0001-01-01 00:00:00 +0000 UTC }
	// true
	// &{test2 status2 group/app2 0001-01-01 00:00:00 +0000 UTC }
	// true
	// &{test3 status3 app3 0001-01-01 00:00:00 +0000 UTC }
	// false
}
