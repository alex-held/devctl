package plugins

import (
	"fmt"
	"os"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"

	"github.com/alex-held/devctl-kit/pkg/devctlpath"
)

type Store interface {
	Register(kind Kind, string string) error
	List(kind Kind) (plugins map[string]string, err error)
}

type pluginsFile struct {
	SDK map[string]string `yaml:"sdk"`
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
		path := s.Pather.ConfigRoot("plugins", name+".so")
		if file.SDK == nil {
			file.SDK = map[string]string{}
		}
		file.SDK[name] = path
	default:
		return fmt.Errorf("failed to register plugin with unsupported pluginkind '%v'", kind)
	}
	return s.save(file)
}

func (s *store) List(kind Kind) (plugins map[string]string, err error) {
	file, err := s.load()
	if err != nil {
		return nil, err
	}
	switch kind {
	case SDK:
		return file.SDK, err
	default:
		return nil, fmt.Errorf("failed to list registered plugins with unsupported pluginkind '%v'", kind)
	}
}

func (s *store) filepath() string { return s.Pather.ConfigRoot("plugins.yaml") }

func (s *store) load() (file *pluginsFile, err error) {
	path := s.filepath()
	file = &pluginsFile{SDK: map[string]string{}}

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
	return afero.WriteFile(s.Fs, path, bytes, os.ModePerm)
}
