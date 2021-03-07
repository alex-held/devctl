package testutils

import (
	"fmt"
	"sync"
	"time"

	"github.com/franela/goblin"
	_ "github.com/franela/goblin"
)

func (f *DetailedWithHooksReporter) CurrentRunID() int {
	return f.ends + f.failures + f.pass + f.fails
}

type DetailedWithHooksReporter struct {
	Next             func()
	DetailedReporter goblin.DetailedReporter
	Describes        []string
	Fails            []string
	Passes           []string
	Pending          []string
	Excluded         []string
	Failures         []*goblin.Failure

	ends, failures, pass, fails, describes, excluded, pending int
	executionTime                                             time.Duration
	totalExecutionTime                                        time.Duration
	executionTimeMu                                           sync.RWMutex
	beginFlag, endFlag                                        bool
}

func (r *DetailedWithHooksReporter) Failure(failure *goblin.Failure) {
	fmt.Printf("\n\n\n ---------------- Failure -------------------- \n\n\n")

	r.failures++
	r.Failures = append(r.Failures, failure)
}

func (r *DetailedWithHooksReporter) BeginDescribe(name string) {
	r.Describes = append(r.Describes, name)
	r.describes++
	fmt.Printf("\n\n\n ---------------- BeginDescribe -------------------- \n\n\n")
	fmt.Printf("\n\n\n ---------------- %s -------------------- \n\n\n", name)
	fmt.Printf("\n\n\n ---------------------------------------- \n\n\n")
}

func (r *DetailedWithHooksReporter) EndDescribe() {
	fmt.Printf("\n\n\n ---------------- EndDescribe -------------------- \n\n\n")
	r.ends++
}

func (r *DetailedWithHooksReporter) ItFailed(name string) {
	fmt.Printf("\n\n\n ---------------- ItFails -------------------- \n\n\n")
	r.Fails = append(r.Fails, name)
	r.fails++
}

func (r *DetailedWithHooksReporter) ItPassed(name string) {
	fmt.Printf("\n\n\n ---------------- ItPassed -------------------- \n\n\n")
	r.Passes = append(r.Passes, name)
	r.pass++
}

func (r *DetailedWithHooksReporter) ItIsPending(name string) {
	fmt.Printf("\n\n\n ---------------- ItIsPending -------------------- \n\n\n")
	r.Pending = append(r.Pending, name)
	r.pending++
}

func (r *DetailedWithHooksReporter) ItIsExcluded(name string) {
	fmt.Printf("\n\n\n ---------------- ItIsExcluded -------------------- \n\n\n")
	r.Excluded = append(r.Excluded, name)
	r.excluded++
}

func (r *DetailedWithHooksReporter) ItTook(duration time.Duration) {
	fmt.Printf("\n\n\n ---------------- ItTook -------------------- \n\n\n")

	fmt.Printf("\n\n\n IT TOOK ---- DONE \n\n\n")
	r.executionTimeMu.Lock()
	defer r.executionTimeMu.Unlock()
	r.executionTime = duration
	r.totalExecutionTime += duration
}

func (r *DetailedWithHooksReporter) Begin() {
	fmt.Printf("\n\n\n ---------------- BEGIN -------------------- \n\n\n")
	r.beginFlag = true
}

func (r *DetailedWithHooksReporter) End() {
	fmt.Printf("\n\n\n ---------------- END -------------------- \n\n\n")
	r.endFlag = true
}
