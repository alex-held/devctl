package plugins

import (
	"context"
	"fmt"
	"io"
	"os"
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
