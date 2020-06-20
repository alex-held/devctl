package kernal

import (
	"fmt"

	"github.com/alex-held/dev-env/config"
)

type gitInstaller struct {
	pathFactory config.PathFactory
	started     bool
	finished    bool
	output      chan string
	error       error
}

func (g *gitInstaller) Started() (finished bool)             { return g.started }
func (g *gitInstaller) Output() *chan string                 { return &g.output }
func (g *gitInstaller) Finished() (finished bool, err error) { return g.finished, g.error }

func NewGitInstaller(pathFactory config.PathFactory) Installer {
	return &gitInstaller{
		pathFactory: pathFactory,
		started:     false,
		finished:    false,
		output:      make(chan string, 10),
	}
}

func (g *gitInstaller) finish() {
	g.finished = true
	close(g.output)
}

func (g *gitInstaller) Install(spec config.Spec) {
	defer g.finish()
	g.started = true

	directory := g.pathFactory.GetPkgDir(spec.Name, spec.Version)
	g.output <- fmt.Sprintf("Cloning %s into %s", spec.Repo, directory)

	if spec.Type != "git" {
		err := fmt.Errorf("Spec of type %s cannot be installed by gitInstaller", spec.Type)
		g.error = err
	} else {
		g.output <- fmt.Sprintf("Installing %s", spec.Name)
	}
}

func (g *gitInstaller) Uninstall(spec config.Spec) {
	defer g.finish()
	g.started = true
	directory := g.pathFactory.GetPkgDir(spec.Name, spec.Version)
	g.output <- fmt.Sprintf("Uninstalling %s from directory %s\n", spec.Name, directory)
	if spec.Type != "git" {
		g.error = fmt.Errorf("Spec of type %s cannot be uninstalled by gitInstaller\n", spec.Type)
	}
}
