package execution

import (
	"fmt"
	. "github.com/alex-held/dev-env/manifest"
	assert2 "github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var execCommand = exec.Command

func RunDevEnvCommand(command DevEnvCommand) ([]byte, error) {
	cmd := exec.Command(command.Command, command.Args...)
	return cmd.CombinedOutput()
}

func RunDummy(cmd *exec.Cmd) ([]byte, error) {
	out, err := cmd.CombinedOutput()
	return out, err
}

func TestName(t *testing.T) {
	execCommand := fakeExecCommand("echo", "alex")
	out, err := RunDummy(execCommand)

	assert2.NoError(t, err)
	assert2.Equal(t, "ALEX", string(out))
}

func TestDevEnvCommand(t *testing.T) {

	command := DevEnvCommand{
		Command: "echo",
		Args:    []string{"alex"},
	}

	out, err := RunDevEnvCommand(command)

	assert2.NoError(t, err)
	assert2.Equal(t, "ALEX", string(out))
}

type TestCommandFactory struct {
	t *testing.T
}

func (factory *TestCommandFactory) CreateDevEnv(command DevEnvCommand) (cmd *exec.Cmd) {
	name := factory.t.Name()
	fmt.Println("TEST-NAME: " + name)

	return CreateCommand(command)
}

func NewTestCommandExecutor(t *testing.T) *CommandExecutor {
	executor := NewCommandExecutor(&Manifest{}, func(str string) {
		_, _ = os.Stdout.WriteString(fmt.Sprint(str))
	})
	executor.CommandFactory = &TestCommandFactory{t: t}
	return executor
}

func TestDevEnvCommandWithFactory(t *testing.T) {

	executor := NewTestCommandExecutor(t)
	command := DevEnvCommand{
		Command: "echo",
		Args:    []string{"alex"},
	}

	out, err := executor.executeDevEnvCommand(command)

	assert2.NoError(t, err)
	assert2.Equal(t, "ALEX", out)
}

func CreateCommand(command DevEnvCommand) (cmd *exec.Cmd) {
	s := []string{command.Command}
	s = append(s, command.Args...)
	cs := []string{"-test.run=TestHelperProcess", "--"}
	cs = append(cs, s...)
	cmd = exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func fakeExecCommand(s ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--"}
	cs = append(cs, s...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)
	args := os.Args

	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}

	if len(args) == 0 {
		_, _ = fmt.Fprintf(os.Stderr, "No command\n")
		os.Exit(2)
	}

	cmd, args := args[0], args[1:]
	switch cmd {
	case "mkdir":
		for _, s := range args {
			fmt.Print("MKDIR")
			fmt.Println(s)
		}
	case "echo":
		fmt.Print(strings.ToUpper(args[0]))
	default:
		_, _ = fmt.Fprintf(os.Stderr, "Unknown command %q\n", cmd)
		os.Exit(2)
	}

	os.Exit(0)
}

/*
// TestHelperProcess isn't a real test. It's used as a helper process
// for TestParameterRun.
func TestHelperProcess(*testing.T) {
    if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
        return
    }
    defer os.Exit(0)

    args := os.Args
    for len(args) > 0 {
        if args[0] == "--" {
            args = args[1:]
            break
        }
        args = args[1:]
    }
    if len(args) == 0 {
        fmt.Fprintf(os.Stderr, "No command\n")
        os.Exit(2)
    }

    cmd, args := args[0], args[1:]
    switch cmd {
    case "echo":
        iargs := []interface{}{}
        for _, s := range args {
            iargs = append(iargs, s)
        }
        fmt.Println(iargs...)
    case "echoenv":
        for _, s := range args {
            fmt.Println(os.Getenv(s))
        }
        os.Exit(0)
    case "cat":
        if len(args) == 0 {
            io.Copy(os.Stdout, os.Stdin)
            return
        }
        exit := 0
        for _, fn := range args {
            f, err := os.Open(fn)
            if err != nil {
                fmt.Fprintf(os.Stderr, "Error: %v\n", err)
                exit = 2
            } else {
                defer f.Close()
                io.Copy(os.Stdout, f)
            }
        }
        os.Exit(exit)
    case "pipetest":
        bufr := bufio.NewReader(os.Stdin)
        for {
            line, _, err := bufr.ReadLine()
            if err == io.EOF {
                break
            } else if err != nil {
                os.Exit(1)
            }
            if bytes.HasPrefix(line, []byte("O:")) {
                os.Stdout.Write(line)
                os.Stdout.Write([]byte{'\n'})
            } else if bytes.HasPrefix(line, []byte("E:")) {
                os.Stderr.Write(line)
                os.Stderr.Write([]byte{'\n'})
            } else {
                os.Exit(1)
            }
        }
    case "stdinClose":
        b, err := ioutil.ReadAll(os.Stdin)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            os.Exit(1)
        }
        if s := string(b); s != stdinCloseTestString {
            fmt.Fprintf(os.Stderr, "Error: Read %q, want %q", s, stdinCloseTestString)
            os.Exit(1)
        }
        os.Exit(0)
    case "exit":
        n, _ := strconv.Atoi(args[0])
        os.Exit(n)
    case "describefiles":
        f := os.NewFile(3, fmt.Sprintf("fd3"))
        ln, err := net.FileListener(f)
        if err == nil {
            fmt.Printf("fd3: listener %s\n", ln.Addr())
            ln.Close()
        }
        os.Exit(0)
    case "extraFilesAndPipes":
        n, _ := strconv.Atoi(args[0])
        pipes := make([]*os.File, n)
        for i := 0; i < n; i++ {
            pipes[i] = os.NewFile(uintptr(3+i), strconv.Itoa(i))
        }
        response := ""
        for i, r := range pipes {
            ch := make(chan string, 1)
            go func(c chan string) {
                buf := make([]byte, 10)
                n, err := r.Read(buf)
                if err != nil {
                    fmt.Fprintf(os.Stderr, "Child: read error: %v on pipe %d\n", err, i)
                    os.Exit(1)
                }
                c <- string(buf[:n])
                close(c)
            }(ch)
            select {
            case m := <-ch:
                response = response + m
            case <-time.After(5 * time.Second):
                fmt.Fprintf(os.Stderr, "Child: Timeout reading from pipe: %d\n", i)
                os.Exit(1)
            }
        }
        fmt.Fprintf(os.Stderr, "child: %s", response)
        os.Exit(0)
    case "exec":
        cmd := exec.Command(args[1])
        cmd.Dir = args[0]
        output, err := cmd.CombinedOutput()
        if err != nil {
            fmt.Fprintf(os.Stderr, "Child: %s %s", err, string(output))
            os.Exit(1)
        }
        fmt.Printf("%s", string(output))
        os.Exit(0)
    case "lookpath":
        p, err := exec.LookPath(args[0])
        if err != nil {
            fmt.Fprintf(os.Stderr, "LookPath failed: %v\n", err)
            os.Exit(1)
        }
        fmt.Print(p)
        os.Exit(0)
    case "stderrfail":
        fmt.Fprintf(os.Stderr, "some stderr text\n")
        os.Exit(1)
    case "sleep":
        time.Sleep(3 * time.Second)
        os.Exit(0)
    default:
        fmt.Fprintf(os.Stderr, "Unknown command %q\n", cmd)
        os.Exit(2)
    }
}
*/
