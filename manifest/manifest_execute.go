package manifest

/*
import (
    "fmt"
    . "github.com/alex-held/dev-env/execution"
    "github.com/spf13/afero"
    "os"
    "os/exec"
    "path"
)

func (cmd LinkCommand) execute(executor CommandExecutor) error {

    targetParentDirectory := path.Dir(cmd.Link.Target)

    mkdir := exec.Command("mkdir", []string{"-p", targetParentDirectory}...)
    ln := exec.Command("ln", []string{"-s", cmd.Link.Source, cmd.Link.Target}...)

    formatAppend(executor.Appender, "%-10s\n", "[Link]")

    exists, err := afero.Exists(afero.NewOsFs(), targetParentDirectory)
    if !exists {

        formatAppend(executor.Appender, "Creating target directory %s\n", targetParentDirectory)

        if !executor.Options.DryRun {
            err = mkdir.Start()
        }

        if err != nil {
            return fmt.Errorf("Error creating linking target directory %s", targetParentDirectory)
        }
    }

    formatAppend(executor.Appender, "Creating link %s -> %s\n", cmd.Link.Source, cmd.Link.Target)

    if !executor.Options.DryRun {
        err = ln.Start()
    }
    if err != nil {
        return fmt.Errorf("Error linking target  %s -> %s\n", cmd.Link.Source, cmd.Link.Target)
    }

    formatAppend(executor.Appender, "\n")
    return nil
}


func (cmd DevEnvCommand) execute(appender func(str string)) error {
    command := exec.Command(cmd.Command, cmd.Args...)

    formatted := cmd.Format()
    formatAppend(appender, "[Command]\n")
    formatAppend(appender, "Executing command: '%s'\n", formatted)
    err := command.Start()

    if err != nil {
        formatAppend(appender, "\n")
        return err
    }

    formatAppend(appender, "\n")
    return nil
}
func (pipe Pipe) execute(appender func(str string)) error {
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

    formatAppend(appender, "[Pipe]\n")
    formatAppend(appender, "\n")

    for i, command := range orderedCommands {
        formatted := pipe.Commands[i].Format()
        formatAppend(appender, "Executing %d/%d %s\n", i, len(orderedCommands), formatted)

        err := command.Start()

        if err != nil {
            formatAppend(appender, "Error '%s' executing pipe command '%s'", err.Error(), formatted)
            return err
        }
    }

    formatAppend(appender, "\n")
    return nil
}*/
