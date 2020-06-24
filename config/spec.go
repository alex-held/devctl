package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

type SpecPackage struct {
	Name        string
	Version     string
	Tags        []string `yaml:"tags,flow"`
	Repo        string
	Description string
}

type Spec struct {
	Package   SpecPackage       `yaml:"package"`
	Variables map[string]string `yaml:"variables,omitempty"`
	Install   struct {
		Instructions []string
	} `yaml:"install,omitempty"`
	UninstallCmds []string `yaml:"uninstall"`
}

type SpecFile struct {
	Spec
	Path string
}

func (spec *Spec) GetInstallInstructions(path PathFactory) (instructions []string, err error) {
	vars := specVars{
		DevEnvHome:      path.GetDevEnvHome(),
		Package:         spec.Package.Name,
		Version:         spec.Package.Version,
		InstallLocation: path.GetPkgDir(spec.Package.Name, spec.Package.Version),
	}
	for _, cmd := range spec.Install.Instructions {
		for key, value := range spec.Variables {
			cmd = strings.ReplaceAll(cmd, "{{ $"+key+" }}", value)
		}
		tmp, err := template.New("cmd").Parse(cmd)
		if err != nil {
			return nil, err
		}
		sb := new(strings.Builder)
		err = tmp.Execute(sb, vars)
		if err != nil {
			return nil, err
		}
		instructions = append(instructions, sb.String())
	}
	return instructions, nil
}

type specVars struct {
	DevEnvHome      string
	Package         string
	Version         string
	InstallLocation string
}

func FromFile(path string) (specFile SpecFile, err error) {
	specFile = SpecFile{Path: path}
	fd, err := os.Open(path)
	if err != nil {
		return specFile, err
	}
	defer fd.Close()
	d := yaml.NewDecoder(fd)
	d.KnownFields(true)

	err = d.Decode(&specFile.Spec)
	return specFile, err
}

func ReadDir(path string) ([]SpecFile, error) {
	log.Println("Reading " + path)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("could not read %s: %w", path, err)
	}
	packages := make([]SpecFile, 0, len(files))
	for _, f := range files {
		// skip directories
		if f.IsDir() {
			continue
		}
		name := filepath.Join(path, f.Name())
		specFile, err := FromFile(name)
		if err != nil {
			fmt.Printf("Warning: could not read %s: %v", name, err)
			continue
		}
		packages = append(packages, specFile)
	}

	return packages, nil
}
