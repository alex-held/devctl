package queue

import "testing"

func TestTaskQueuer(t *testing.T) {





	
}






/*
import (
	"fmt"
	"sync"
	"testing"
)

type queue struct {
	jobs []Job
}

func (q *queue) Enqueue(j Job) {
	q.jobs = append(q.jobs, j)
}

type Jobber interface {
	Job()
}

type Job interface {
	Name() string
	Error() error
}

type Requirer interface {
	Required() bool
}

type Depender interface {
	Register(dependencyResolverC <-chan interface{})
	Wait() <-chan interface{}
}

type job struct {
	name     string
	v        interface{}
	callback func(interface{})
}

func (j *job) Job() {
	j.callback(j.v)
}

type dependendJob struct {
	dependsOnSuccess []chan interface{}
	dependsOn        []chan interface{}
	wait             chan interface{}
	value            interface{}
}

func Test1(t *testing.T) {

	var (
		result interface{}
		wg     sync.WaitGroup
	)

	_ = &queue{}

	job1 := &job{
		name: "Job 1",
		v:    1,
		callback: func(v interface{}) {
			result = fmt.Sprintf("%d_bar", v)
			result = v
			wg.Done()
		},
	}

	go job1.Job()
	wg.Wait()

}
*/
