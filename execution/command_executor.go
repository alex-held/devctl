package execution

import (
	"fmt"
	. "github.com/alex-held/dev-env/manifest"
	"github.com/spf13/afero"
	"os"
	"os/exec"
	"path"
)

type DefaultCommandFactory struct {
}

func (d DefaultCommandFactory) CreateDevEnv(command DevEnvCommand) (cmd *exec.Cmd) {
	cmd = exec.Command(command.Command, command.Args...)
	return cmd
}

type CommandFactory interface {
	CreateDevEnv(command DevEnvCommand) (cmd *exec.Cmd)
}

type CommandExecutor struct {
	Manifest       *Manifest
	Appender       func(str string)
	FS             *afero.Fs
	Options        *CommandExecutorOptions
	CommandFactory CommandFactory
}

type CommandExecutorOptions struct {
	DryRun bool
}

func NewCommandExecutorOptions() *CommandExecutorOptions {
	return &CommandExecutorOptions{
		DryRun: false,
	}
}

func (ce *CommandExecutor) DirectoryExist(path string) bool {
	exists, _ := afero.Exists(*ce.FS, path)
	return exists
}

func NewCommandExecutor(manifest *Manifest, appender func(str string)) *CommandExecutor {
	options := NewCommandExecutorOptions()
	fs := afero.NewOsFs()
	return &CommandExecutor{
		Manifest:       manifest,
		Appender:       appender,
		FS:             &fs,
		Options:        options,
		CommandFactory: DefaultCommandFactory{},
	}
}

func (ce *CommandExecutor) Execute() (output string, err error) {
	instructions := ce.Manifest.ResolveInstructions()
	var out []byte
	for _, instruction := range instructions {

		switch executable := instruction.(type) {
		case DevEnvCommand:
			output, err = ce.executeDevEnvCommand(executable)
		case LinkCommand:
			err = ce.executeLinkCommand(executable)
		case Pipe:
			err = ce.executePipe(executable)
		}
		if err != nil {
			fmt.Printf("Error executing instruction %s\n%+v", err.Error(), instruction)
			return "", err
		}
	}
	return string(out), nil
}

func (ce *CommandExecutor) append(format string, args ...interface{}) {
	formatted := fmt.Sprintf(format, args...)
	ce.Appender(formatted)
}

func (ce *CommandExecutor) executePipe(pipe Pipe) error {

	cmd1 := exec.Command(pipe.Commands[0].Command, pipe.Commands[0].Args...)

	orderedCommands := []*exec.Cmd{cmd1}
	for i, command := range pipe.Commands {

		if i == 0 {
			continue
		}

		cNext := exec.Command(command.Command, command.Args...)
		cNext.Stdin, _ = cmd1.StdoutPipe()
		cNext.Stdout = os.Stdout
		orderedCommands = append(orderedCommands, cNext)
	}

	ce.append("[Pipe]\n")
	ce.append("\n")

	for i, command := range orderedCommands {
		formatted := pipe.Commands[i].Format()
		ce.append("Executing %d/%d %s\n", i, len(orderedCommands), formatted)

		if !ce.Options.DryRun {
			err := command.Start()
			if err != nil {
				ce.append("Error '%s' executing pipe command '%s'", err.Error(), formatted)
				ce.append("\n")
				return err
			}
		}

	}

	ce.append("\n")
	return nil
}

func (ce *CommandExecutor) executeDevEnvCommand(executable DevEnvCommand) (output string, err error) {

	command := ce.CommandFactory.CreateDevEnv(executable)
	var out []byte
	formatted := executable.Format()
	ce.append("[Command]\n")
	ce.append("Executing command: '%s'\n", formatted)

	if !ce.Options.DryRun {
		out, err = command.CombinedOutput()
		if err != nil {
			ce.append("\n")
			return "", err
		}
	}

	ce.append("\n")
	return string(out), nil
}

func (ce *CommandExecutor) executeLinkCommand(cmd LinkCommand) error {
	targetParentDirectory := path.Dir(cmd.Link.Target)

	mkdir := exec.Command("mkdir", []string{"-p", targetParentDirectory}...)
	ln := exec.Command("ln", []string{"-s", cmd.Link.Source, cmd.Link.Target}...)
	var err error
	ce.append("%-10s\n", "[Link]")

	if !ce.DirectoryExist(targetParentDirectory) {

		ce.append("Creating target directory %s\n", targetParentDirectory)

		if !ce.Options.DryRun {
			err = mkdir.Start()
		}

		if err != nil {
			return fmt.Errorf("Error creating linking target directory %s", targetParentDirectory)
		}
	}

	ce.append("Creating link %s -> %s\n", cmd.Link.Source, cmd.Link.Target)

	if !ce.Options.DryRun {
		err = ln.Start()
	}
	if err != nil {
		return fmt.Errorf("Error linking target  %s -> %s\n", cmd.Link.Source, cmd.Link.Target)
	}

	ce.append("\n")
	return nil
}
