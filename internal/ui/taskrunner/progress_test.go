package taskrunner

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"
)

type NoOpPlugin struct {
	Out io.Writer
}

type TestNoOpPlugin struct {
	NoOpPlugin
	error
}

func (t TestNoOpPlugin) ExecuteCommand(_ context.Context, root string, args []string) error {
	fmt.Fprintf(t.Out, "Executing NoOpPlugin..\tRoot=%s\tArgs=%v\n", root, args)
	if t.error != nil {
		fmt.Fprintf(t.Out, "[NoOpPlugin] Error occurred! ErrorF=%+v\n", t.error)
	}
	fmt.Fprintln(t.Out, "[NoOpPlugin] Success! occurred!")
	return t.error
}

func (NoOpPlugin) PluginName() string { return "NoOpPlugin" }

func (p NoOpPlugin) ExecuteCommand(_ context.Context, root string, args []string) error {
	out := p.Out
	if out == nil {
		out = os.Stdout
	}
	fmt.Fprintf(out, "Executing NoOpPlugin..\tRoot=%s\tArgs=%v\n", root, args)
	return nil
}

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
		Plugin: NoOpPlugin{
			Out: io.Discard,
		},
		Description: "Downloading go sdk",
		Root:        "test",
		Args:        []string{},
	},
	{
		Plugin: NoOpPlugin{
			Out: io.Discard,
		},
		Description: "Extracting go sdk",
		Root:        "test",
		Args:        nil,
	},
	{
		Plugin: NoOpPlugin{
			Out: io.Discard,
		},
		Description: "Linking go sdk",
		Root:        "test",
		Args:        nil,
	},
	{
		Plugin: NoOpPlugin{
			Out: io.Discard,
		},
		Description: "Listing go sdk",
		Root:        "test",
		Args:        nil,
	},
	/*{
		Plugin: NoOpPlugin{
			Out: io.Discard,
		},
		Description: "NoOp Task 5",
		Root:        "test",
		Args:        nil,
	},*/
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
