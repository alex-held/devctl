package plugins

import (
	"errors"
	"plugin"
)

type pluginRegistry struct {
	SDKPlugins []SDKPlugin
}

// Plugin is a plugin loaded from a file
type DevCtlPlugin struct {
	// Name of the plugin e.g rabbitmq
	Name string
	// Type of the plugin e.g broker
	Type string
	// Path specifies the import path
	Path string
	// NewFunc creates an instance of the plugin
	NewFunc interface{}
}

func (pr *pluginRegistry) Load(p string) (*plugin.Plugin, error) {
	plug, err := plugin.Open(p)
	if err != nil {
		return nil, err
	}
	s, err := plug.Lookup("Plugin")
	if err != nil {
		return nil, err
	}
	pl, ok := s.(*plugin.Plugin)
	if !ok {
		return nil, errors.New("could not cast Plugin object")
	}
	return pl, nil
}
