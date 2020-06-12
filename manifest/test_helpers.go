package manifest

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strconv"
    "strings"
    "testing"
    _ "github.com/ganbarodigital/go_pipe/v5"
    scriptish "github.com/ganbarodigital/go_scriptish"
    "github.com/spf13/afero"
)

type CommandType = int

const (
    Single CommandType = iota
    List
    Piped
)

type CommandExecutionManager interface {
    Commandource
    GetFs() afero.Fs
    GetFactory() (factory *CommandFactory)
    Write(buffer *bytes.Buffer) (err error)
    Execute() (out []byte, err error)
}

type CommandFactory interface {
    Create(command string, args ...string) (cmd *exec.Cmd)
}

var execCommand = exec.Command

func fakeExecCommand(command string, args ...string) *exec.Cmd {
    cs := []string{"-test.run=TestHelperProcess", "--", command}
    cs = append(cs, args...)
    cmd := exec.Command(os.Args[0], cs...)
    cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
    return cmd
}

type testCmdFactory struct {
    create func(command string, args ...string) (cmd *exec.Cmd)
}

func (t testCmdFactory) Create(command string, args ...string) (cmd *exec.Cmd) {
    return t.create(command, args...)
}

func NewTestCommandFactory() CommandFactory {
    var factory CommandFactory = testCmdFactory{create: fakeExecCommand}
    return factory
}

func FakeExecCommand(command string, args ...string) *exec.Cmd {
    cs := []string{"-test.run=TestExecCommandHelper", "--", command}
    cs = append(cs, args...)
    cmd := exec.Command(os.Args[0], cs...)
    cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1",
        "STDOUT=" + command,
        "EXIT_STATUS=" + "echo hallo welt"}
    return cmd
}

func TestExecCommandHelper(t *testing.T) {
    if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
        return
    }

    // println("Mocked stdout:", os.Getenv("STDOUT"))
    fmt.Fprintf(os.Stdout, os.Getenv("STDOUT"))
    i, _ := strconv.Atoi(os.Getenv("EXIT_STATUS"))
    os.Exit(i)
}

type Command interface {
    GetPipeline() PipelineCommand
    WriteTo(writer bufio.Writer) (err error)
    Execute(writer bufio.Writer) (err error)
    GetCommandType() CommandType
    Print()
    //resolve(variables map[string]string) Command
}

type DefaultCommand struct {
    *CommandType
    Command string
    Args    []string
}

func (cmd *DefaultCommand) GetCommandType() CommandType {
    return *cmd.CommandType
}

type DevEnvCommand struct {
    Command string   `json:"command,omitempty"`
    Args    []string `json:"args,omitempty"`
}

func (d DevEnvCommand) GetPipeline() (cmd PipelineCommand) {
    cmd = PipelineCommand{
        scriptish.NewPipeline(scriptish.Exec(d.Args...)),
        Single,
    }
    return cmd
}


type Link struct {
    Source string `json:"source"`
    Target string `json:"target"`
}


type Pipe struct {
    Commands []DevEnvCommand `json:"commands, omitempty"`
    IsPiping bool
}

func (p Pipe) GetPipeline() PipelineCommand {
    panic("implement me")
}


func (p Pipe) Execute(writer bufio.Writer) (err error) {
    panic("implement me")
}

func (p Pipe) GetCommandType() CommandType {
    return Piped
}

type LinkCommand struct {
    Link Link
}

func (d LinkCommand) Execute(writer bufio.Writer) (err error) {
    panic("implement me")
}

func (d LinkCommand) GetCommandType() CommandType {
   return Single
}


func (d DevEnvCommand) Execute(writer bufio.Writer) (err error) {
    panic("implement me")
}

func (d DevEnvCommand) GetCommandType() CommandType { return Single }

type PipelineCommand struct {
    Commands *scriptish.Sequence
    CommandType
}

func (p PipelineCommand) Print() {
    writer := bufio.NewWriter(os.Stdout)
    _ = p.WriteTo(*writer)
}

func (p PipelineCommand) GetPipeline() PipelineCommand {
    return p
}

func (p PipelineCommand) WriteTo(bufio.Writer) (err error) {
    panic("implement me")
}

func (p PipelineCommand) Execute(writer bufio.Writer) (err error) {
   return p.Execute(writer)
}

func (p PipelineCommand) GetCommandType() CommandType {
    panic("implement me")
}

func NewPipeline(cmdType CommandType, steps ...scriptish.Command) *PipelineCommand {
    var pipeline = scriptish.NewSequence(
        scriptish.Exec("", ""),
        scriptish.Exec("", ""),
    )

    pCmd := PipelineCommand{
        Commands: pipeline,
        CommandType: cmdType,
    }
    return &pCmd
}

func formatInternal(cmd PipelineCommand) (format string) {
    seq := cmd.Commands
    format, _ = seq.TrimmedString()
    return format
}

type Commandource interface {
    GetInstructions() []Command
}


type CommandExecutor struct {
    Factory CommandFactory
    Testing *testing.T
    FS      *afero.Fs
    Source  Commandource
    Writer  strings.Builder
    Options CommandExecutorOptions
}

func (manager CommandExecutor) GetInstructions() []Command {
   return manager.Source.GetInstructions()
}

func (manager CommandExecutor) GetFs() afero.Fs {
   return *manager.FS
}

func (c CommandExecutor) Write(buffer *bytes.Buffer) (err error) {
w := bufio.NewWriter(buffer)
w.WriteString("func (c *CommandExecutor) Write(writer bufio.Writer) (err error) {")
if err != nil {
err = w.Flush()
}
_, err = w.WriteString("func (c *CommandExecutor) Write(writer bufio.Writer) (err error)")
return err
}

func (c CommandExecutor) GetFactory() (factory *CommandFactory) {
    return &c.Factory
}

type CommandExecutorOptions struct {
    DryRun bool
}
