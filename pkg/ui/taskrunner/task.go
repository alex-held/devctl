package taskrunner

import (
	"context"
	"fmt"

	plugins2 "github.com/alex-held/devctl/pkg/plugins"
)

type SimpleTask struct {
	Description string
	Action      func(ctx context.Context) error
}
func (t *SimpleTask) Describe() string { return t.Description }
func (t *SimpleTask) Task(ctx context.Context) (err error) { return t.Action(ctx) }

type Task struct {
	Description string
	Root        string
	Plugin      plugins2.Executor
	Args        []string
}

func (t *Task) Describe() string {
	return fmt.Sprintf("Title:\t%s\tPlugin:\t%s\tArgs:\t%v\n",
		t.Description,
		t.Plugin.PluginName(),
		t.Args,
	)
}

func (t *Task) Task(ctx context.Context) (err error) {
	return t.Plugin.ExecuteCommand(ctx, t.Root, t.Args)
}

type Tasks []Tasker

func (t Tasks) Filter(filterFn func(tasker Tasker) bool) (result Tasks) {
	for _, tasker := range t {
		if filterFn(tasker) {
			result = append(result, tasker)
		}
	}
	return result
}

func (t Tasks) MapToString(mapFn func(tasker Tasker) string) (result []string) {
	for _, tasker := range t {
		mapped := mapFn(tasker)
		result = append(result, mapped)
	}
	return result
}

func NewConditionalTask(desc string, task TaskActionFn, shouldExecuteFn ConditionalExecutorFn) Tasker {
	return &ConditionalTask{
		Description:   desc,
		Action:        task,
		ShouldExecute: shouldExecuteFn,
	}
}

type ConditionalTask struct {
	Description   string
	Action        func(ctx context.Context) error
	ShouldExecute func() bool
}

func (t *ConditionalTask) Describe() string {
	return t.Description
}

func (t *ConditionalTask) Task(ctx context.Context) (err error) {
	if t.ShouldExecute() {
		return t.Action(ctx)
	}
	return nil
}
