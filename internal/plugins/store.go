package plugins

import (
	"fmt"
	"io/fs"

	"github.com/alex-held/devctl/pkg/devctlpath"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

type Kind int

const (
	SDK Kind = iota
)

type Store interface {
	Register(kind Kind, string string) error
	List(kind Kind) (plugins []string, err error)
}

type pluginsFile struct {
	SDK []string `yaml:"sdk"`
}

type store struct {
	Pather devctlpath.Pather
	Fs     afero.Fs
}

func (s *store) Register(kind Kind, name string) (err error) {
	file, err := s.load()
	if err != nil {
		return err
	}
	switch kind {
	case SDK:
		file.SDK = append(file.SDK, name)
	default:
		return fmt.Errorf("failed to register plugin with unsupported pluginkind '%v'", kind)
	}
	return s.save(file)
}

func (s *store) List(kind Kind) (plugins []string, err error) {
	file, err := s.load()
	if err != nil {
		return nil, err
	}
	switch kind {
	case SDK:
		return file.SDK, nil
	default:
		return nil, fmt.Errorf("failed to list registered plugins with unsupported pluginkind '%v'", kind)
	}
}

func (s *store) filepath() string { return s.Pather.ConfigRoot("plugins.yaml") }

func (s *store) load() (file *pluginsFile, err error) {
	path := s.filepath()
	file = &pluginsFile{SDK: []string{}}

	if exist, _ := afero.Exists(s.Fs, path); !exist {
		return file, nil
	}
	bytes, err := afero.ReadFile(s.Fs, path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(bytes, file)
	return file, err
}

func (s *store) save(file *pluginsFile) (err error) {
	path := s.filepath()
	bytes, err := yaml.Marshal(file)
	if err != nil {
		return err
	}
	return afero.WriteFile(s.Fs, path, bytes, fs.ModePerm)
}
