package taskrunner

import (
	"context"
	"io"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/alex-held/devctl/pkg/plugins"
)

var defaultTestTasks = Tasks{
	&Task{
		Plugin: plugins.NoOpPlugin{
			Out: io.Discard},
		Description: "Downloading go sdk",
		Root:        "test",
		Args:        []string{},
	},
	&Task{
		Plugin: plugins.NoOpPlugin{
			Out: io.Discard,
		},
		Description: "Extracting go sdk",
		Root:        "test",
		Args:        nil,
	},
	&Task{
		Plugin: plugins.NoOpPlugin{
			Out: io.Discard,
		},
		Description: "Linking go sdk",
		Root:        "test",
		Args:        nil,
	},
	&Task{
		Plugin: plugins.NoOpPlugin{
			Out: io.Discard,
		},
		Description: "Listing go sdk",
		Root:        "test",
		Args:        nil,
	},
}

var testDescriber = func(t Tasker) string { return t.Describe() }

func TestNewTaskRunner(t *testing.T) {
	ctx := context.TODO()

	sut := NewTaskRunner(
		WithTitle("TestNewTaskRunner"),
		WithTasks(defaultTestTasks...),
	)

	err := sut.Run(ctx)

	if err != nil {
		t.Fatal(err)
	}
}

func TestNewTaskRunnerWithConditionalTasks(t *testing.T) {
	ctx := context.TODO()
	var tasks Tasks

	tasks = append(tasks, &ConditionalTask{
		Description: "[sdk/go/download] - NOT REQUIRED",
		Action: func(ctx context.Context) error {
			t.Fatalf("Shoukld not execute \"[sdk/go/download] - NOT REQUIRED\"")
			return nil
		},
		ShouldExecute: NotRequired.ShouldExecute,
	},
	)

	sut := NewTaskRunner(
		WithTitle("TestNewTaskRunnerWithConditionalTasks"),
		WithTimeout(500*time.Millisecond),
		WithTasks(tasks...),
	)

	err := sut.Run(ctx)

	if err != nil {
		t.Fatal(err)
	}
}

type Conditional bool

func (c Conditional) ShouldExecute() bool {
	return bool(c)
}

func (c Conditional) Filter(tasks Tasks) Tasks {
	return tasks.Filter(func(tasker Tasker) bool {
		conditionalTask := tasker.(*ConditionalTask)
		return conditionalTask.ShouldExecute() == bool(c)
	})
}

const (
	Required    Conditional = true
	NotRequired Conditional = false
)

var _ = Describe("Tasker", func() {
	var (
		runner Runner
		ctx    context.Context
		desc   GinkgoTestDescription
	)

	Context("ConditionalTask", func() {
		var actualTaskDesc []string

		var tasks = Tasks{
			&ConditionalTask{
				Description: "[sdk/go/installer] - REQUIRED",
				Action: func(ctx context.Context) error {
					actualTaskDesc = append(actualTaskDesc, "[sdk/go/installer] - REQUIRED")
					return nil
				},
				ShouldExecute: Required.ShouldExecute,
			},
			&ConditionalTask{
				Description: "[sdk/go/download] - NOT REQUIRED",
				Action: func(ctx context.Context) error {
					actualTaskDesc = append(actualTaskDesc, "[sdk/go/download] - NOT REQUIRED")
					return nil
				},
				ShouldExecute: NotRequired.ShouldExecute,
			},
		}

		BeforeEach(func() {
			desc = CurrentGinkgoTestDescription()
			ctx = context.TODO()
		})

		When("pipeline contains tasks that are not required", func() {
			BeforeEach(func() {
				runner = NewTaskRunner(
					WithDiscardOutput(),
					WithTitle(desc.TestText),
					WithTasks(NotRequired.Filter(tasks)...),
				)
			})

			It("shouldn't execute the not required tasks", func() {
				Expect(runner.Run(ctx)).To(Succeed())
				Expect(actualTaskDesc).Should(BeEmpty())
			})
		})

		When("pipeline contains tasks that are required", func() {
			BeforeEach(func() {
				runner = NewTaskRunner(
					WithDiscardOutput(),
					WithTitle(desc.TestText),
					WithTasks(Required.Filter(tasks)...),
				)
			})

			It("should execute the required tasks", func() {
				expectedTaskDesc := Required.Filter(tasks).MapToString(testDescriber)
				Expect(runner.Run(ctx)).To(Succeed())
				Expect(actualTaskDesc).Should(Equal(expectedTaskDesc))
			})
		})
	})
})
