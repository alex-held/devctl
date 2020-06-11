package manifest

import (
    "bufio"
    scriptish "github.com/ganbarodigital/go_scriptish"
    "github.com/spf13/afero"
    "os/exec"
    . "path"
    "strings"
)

type DefaultCommandFactory struct {
    afero.Fs
    create func(command string, args ...string) (cmd *exec.Cmd)
}

func (t *DefaultCommandFactory) Create(command string, args ...string) (cmd *exec.Cmd) {
    return t.create(command, args...)
}

func (t *DefaultCommandFactory) CreateCmd(command string, args ...string) (cmd *exec.Cmd) {
    return exec.Command(command, args...)
}

func (t *DefaultCommandFactory) CreateDevEnv(command DevEnvCommand) (cmd *exec.Cmd) {
    cmd = exec.Command(command.Command, command.Args...)
    return cmd
}

func NewCommandExecutorOptions() *CommandExecutorOptions {
    return &CommandExecutorOptions{
        DryRun: false,
    }
}

func NewCommandExecutor(commandoSource Commandource) CommandExecutionManager {
    options := *NewCommandExecutorOptions()
    factory :=NewTestCommandFactory()
    fs := afero.NewOsFs()

    executor := CommandExecutor{
        Factory: factory,
        FS:      &fs,
        Source:  commandoSource,
        Writer:  strings.Builder{},
        Options: options,
    }
    var manager CommandExecutionManager = executor

    return manager
}



func (manager CommandExecutor) Execute() ([]byte, error) {

    var file afero.File
    var instructions = manager.Source.GetInstructions()

    _ = *manager.GetFactory()
    file, _ = (*manager.FS).Create("/Users/dev/temp/execution.log")
    file2, _ := afero.ReadFile(*manager.FS, "/Users/dev/temp/execution.log")
    writer := bufio.NewWriter(file)
    destination := scriptish.NewDest()
    _, _ = writer.WriteString(<-destination.ReadLines() + "\n")

    for _, iCmd := range instructions {
        pl := iCmd.GetPipeline()
        scriptish.WriteToFile("/Users/dev/temp/execution.log")
        _ = pl.Execute(*writer)
    }

   return file2, nil
}

func CreateCommands(manager CommandExecutionManager) (result []PipelineCommand) {
    result = []PipelineCommand{}
    //factory := manager.GetFactory()
    for _, iCommand := range manager.GetInstructions() {
        switch instr := iCommand.(type) {
        case LinkCommand:
            targetParentDirectory := Dir(instr.Link.Target)
            result = append(result, PipelineCommand{
                Commands: scriptish.NewPipeline(
                scriptish.Mkdir(targetParentDirectory, 777),
                scriptish.Exec(instr.Link.Source, instr.Link.Target),
            ),
                CommandType: List})
        case Pipe:
            for _, command := range instr.Commands {
                args := []string{command.Command}
                args = append(args, command.Args...)
                result = append(result, PipelineCommand{
                    Commands:    scriptish.NewPipeline(scriptish.Exec(args...)),
                    CommandType: Piped,
                },
                )
            }
        default:
            //   fmt.Sprintf("%+v", iCommand)
            continue
        }
    }
    return result
}
