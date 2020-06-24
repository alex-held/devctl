package kernal

import (
	"fmt"

	pipes "github.com/ebuchman/go-shell-pipes"

	"github.com/alex-held/dev-env/config"
)

type installer struct {
	pathFactory config.PathFactory
	started     bool
	finished    bool
	output      chan string
	error       error
	options     InstallerOptions
}

func (g *installer) Started() (finished bool)             { return g.started }
func (g *installer) Output() *chan string                 { return &g.output }
func (g *installer) Finished() (finished bool, err error) { return g.finished, g.error }

type InstallerOptions struct {
	dry bool
}

func NewInstaller(pathFactory config.PathFactory, options InstallerOptions) Installer {
	return &installer{
		pathFactory: pathFactory,
		started:     false,
		finished:    false,
		output:      make(chan string, 10),
		options:     options,
	}
}

func (g *installer) finish() {
	g.finished = true
	close(g.output)
}

func (g *installer) Install(spec config.Spec) {
	defer g.finish()
	g.started = true

	directory := g.pathFactory.GetPkgDir(spec.Package.Name, spec.Package.Version)
	g.output <- fmt.Sprintf("Installing %s into directory %s", spec.Package.Name, directory)

	instructions, err := spec.GetInstallInstructions(g.pathFactory)
	installCmdCount := len(instructions)
	if err != nil {
		g.error = err
		return
	}

	for i, line := range instructions {
		g.output <- fmt.Sprintf("[STEP %d/%d]%6s", i+1, installCmdCount, line)
		if !g.options.dry {
			out, err := pipes.RunString(line)
			if err != nil {
				g.error = err
				return
			}
			g.output <- out
		}
	}

}

func (g *installer) Uninstall(spec config.Spec) {
	defer g.finish()
	g.started = true
	directory := g.pathFactory.GetPkgDir(spec.Package.Name, spec.Package.Version)
	g.output <- fmt.Sprintf("Uninstalling %s from directory %s\n", spec.Package.Name, directory)
}
