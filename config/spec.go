package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Spec struct {
	Name          string
	Version       string
	Type          string
	Tags          []string
	Repo          string
	Desc          string
	InstallCmds   []string `yaml:"install"`
	UninstallCmds []string `yaml:"uninstall"`
}

type SpecFile struct {
	Spec
	Path string
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
