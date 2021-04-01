package cli

import (
	"fmt"
	"os/exec"

)

type ExecProxy struct {
	version string
	shims   []ExecProxy
}

type CLIApp struct {
	version string
	shims   map[string]Shim
}

func NewCLIApp(version string) *CLIApp {
	return & CLIApp{
		version: version,
		shims: map[string]Shim{
			"go": &shim{
				name:                 "go",
				version:              "1.0.0",
				nextResponsibleProxy:  nil,
			},
		},
	}
}

func (cli *CLIApp) Version() string {
	return cli.version
}

func (cli *CLIApp) Exec(toolname string, args ...string) (err error) {
	shim, err := cli.GetShimForTool(toolname)
	if err != nil {
		return err
	}
	return shim.Exec(args...)
}

type shim struct {
	name                 string
	version              string
	nextResponsibleProxy *ExecProxy
}

func (s *shim) Version() (string, error) {
	return s.version, nil
}

func (s *shim) Exec(args ...string) error {
	return s.Next().Exec(args...)
}

type Shim interface {
	Version()  (string, error)
	Exec(args ...string) error
}

type CliCommandShim struct {
   Cmd *exec.Cmd
}

func (c *CliCommandShim) Version() (string, error) {
	c.Cmd.Args = append(c.Cmd.Args, "--version" )
	stdout, stderr := c.Cmd.CombinedOutput()
	return string(stdout), stderr
}


func (c *CliCommandShim) Exec(args ...string) error {
	c.Cmd.Args = append(c.Cmd.Args, args...)
	stderr := c.Cmd.Start()
	return stderr
}

func (cli *CLIApp) GetShimForTool(toolname string) (shim Shim, err error) {

	if val, ok := cli.shims[toolname]; ok != false {
		return val, nil
	} else if path, err := exec.LookPath(toolname); err != nil {
		return &CliCommandShim{
			exec.Command(path),
		}, nil
	}
	err = fmt.Errorf("A shim with name %s has not been registered anc could be not  be found on your system.\n", toolname)
	return nil, err
}
