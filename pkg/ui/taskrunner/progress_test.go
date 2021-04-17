package taskrunner

import (
	"context"
	"io"
	"testing"

	"github.com/alex-held/devctl/pkg/plugins"
)

/*
func TestTaskRunner_Run(t *testing.T) {
	ctx := context.TODO()
	sut := taskRunner{
		Title: "TestTaskRunner_Run",
		Tasks: []Task{
			{
				Plugin: NoOpPlugin{
					Out: io.Discard,
				},
				Description: "NoOp Test 1",
				Root:        "test",
				Args:        nil,
			},
			{
				Plugin: NoOpPlugin{
					Out: io.Discard,
				},
				Description: "NoOp Task 2",
				Root:        "test",
				Args:        nil,
			},
			{
				Plugin: NoOpPlugin{
					Out: io.Discard,
				},
				Description: "NoOp Task 3",
				Root:        "test",
				Args:        nil,
			},
			{
				Plugin: NoOpPlugin{
					Out: io.Discard,
				},
				Description: "NoOp Task 4",
				Root:        "test",
				Args:        nil,
			},
			{
				Plugin: NoOpPlugin{
					Out: io.Discard,
				},
				Description: "NoOp Task 5",
				Root:        "test",
				Args:        nil,
			},
		},
	}

	err := sut.Run(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
*/

var defaultTestTasks = Tasks{
	{
		Plugin: plugins.NoOpPlugin{
			Out: io.Discard,
		},
		Description: "Downloading go sdk",
		Root:        "test",
		Args:        []string{},
	},
	{
		Plugin: plugins.NoOpPlugin{
			Out: io.Discard,
		},
		Description: "Extracting go sdk",
		Root:        "test",
		Args:        nil,
	},
	{
		Plugin: plugins.NoOpPlugin{
			Out: io.Discard,
		},
		Description: "Linking go sdk",
		Root:        "test",
		Args:        nil,
	},
	{
		Plugin:plugins. NoOpPlugin{
			Out: io.Discard,
		},
		Description: "Listing go sdk",
		Root:        "test",
		Args:        nil,
	},
}

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
