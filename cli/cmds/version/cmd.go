package version

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/gobuffalo/plugins/plugprint"
	"github.com/spf13/pflag"
)

var _ plugins.Plugin = &Cmd{}
var _ plugio.OutNeeder = &Cmd{}
var _ plugprint.Describer = &Cmd{}
var _ plugprint.FlagPrinter = &Cmd{}

// Cmd is responsible for the `buffalo version` command.
type Cmd struct {
	help   bool
	json   bool
	stdout io.Writer
}

func (cmd *Cmd) SetStdout(w io.Writer) error {
	cmd.stdout = w
	return nil
}

func (cmd *Cmd) PrintFlags(w io.Writer) error {
	flags := cmd.Flags()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}

func (cmd *Cmd) PluginName() string {
	return "version"
}

func (cmd *Cmd) Description() string {
	return "Print the version information"
}

func (cmd Cmd) String() string {
	return cmd.PluginName()
}

func (cmd *Cmd) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(cmd.String(), pflag.ContinueOnError)
	flags.SetOutput(ioutil.Discard)
	flags.BoolVarP(&cmd.help, "help", "h", false, "print this help")
	flags.BoolVarP(&cmd.json, "json", "j", false, "Print information in json format")
	return flags
}

func (cmd *Cmd) Main(ctx context.Context, root string, args []string) error {
	flags := cmd.Flags()
	if err := flags.Parse(args); err != nil {
		return err
	}

	out := cmd.stdout
	if cmd.stdout == nil {
		out = os.Stdout
	}
	if cmd.help {
		return plugprint.Print(out, cmd)
	}

	if !cmd.json {
		fmt.Fprintln(out, Version)
		return nil
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", "    ")
	return enc.Encode(map[string]string{
		"version": Version,
	})
}
