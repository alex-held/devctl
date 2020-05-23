package manifest

import (
	"fmt"
	"sort"
	"strings"
)

type InstructionType int
type Variables []Variable
type Instructions []Step

const (
	Command InstructionType = iota
	CommandPipe
)

func (variables *Variables) ToMap() StringMap {
	result := StringMap{}

	for _, variable := range *variables {
		result[variable.Key] = variable.Value
	}
	return result
}

type StringSliceStringMap map[string]interface{}
type StringMap map[string]string

type Manifest struct {
	Version      string       `json:"version"`
	SDK          string       `json:"sdk"`
	Variables    Variables    `json:"variables,omitempty"`
	Instructions Instructions `json:"instructions,omitempty"`
	Links        []Link       `json:"links,omitempty"`
}

func (m *Manifest) ResolveVariable() Variables {
	var result Variables

	variables := m.populateVariables()

	for key, val := range variables {
		replaced := val
		resolvedVariables := ResolveTemplateValues(val, variables)

		for resolvedKey, value := range resolvedVariables {

			if _, ok := variables[resolvedKey]; !ok {
				variables[resolvedKey] = value
			}
		}

		for vKey, vValue := range variables {
			replaced = strings.ReplaceAll(replaced, vKey, vValue)
		}

		variables[key] = replaced
	}

	for key, value := range variables {
		result = append(result, Variable{
			Key:   key,
			Value: value,
		})
	}

	sort.Sort(result)
	return result
}

func (v Variables) Len() int {
	return len(v)
}

func (v Variables) Swap(i, j int) {
	iVal := v[i]
	jVal := v[j]
	v[i] = jVal
	v[j] = iVal
}

func (v Variables) Less(i, j int) bool {
	return v[i].Key < v[j].Key
}

type Link struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

func NewInstaller(commands ...Instructing) *Installer {
	i := Installer{
		Instructions: map[int]Instructing{},
		Commands:     []DevEnvCommand{},
	}
	for i2, command := range commands {
		i.Instructions[i2] = command
	}
	return &i
}

type Installer struct {
	Instructions map[int]Instructing `json:"instructions"`
	Commands     []DevEnvCommand     `json:"commands"`
}

type Step struct {
	Command *DevEnvCommand  `json:"command,omitempty"`
	Pipe    []DevEnvCommand `json:"pipe,omitempty"`
}

type DevEnvCommand struct {
	Command string   `json:"command,omitempty"`
	Args    []string `json:"args,omitempty"`
}

type Pipe struct {
	Commands []DevEnvCommand `json:"commands, omitempty"`
}

func (step *Step) ToInstruction() Instructing {
	if step.Pipe != nil {
		return Pipe{Commands: step.Pipe}
	}
	if step.Command != nil {
		return DevEnvCommand{
			Command: step.Command.Command,
			Args:    step.Command.Args,
		}
	}
	return nil
}

type Instructing interface {
	Format() string
	Execute() error
	Resolve(variables Variables) Instructing
}

type Variable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (m *Manifest) MustGetVariable(key string) string {
	variable, err := m.GetVariable(key)
	if err != nil {
		panic(err.Error())
	}
	return *variable
}

func (m *Manifest) GetVariable(key string) (*string, error) {
	if key == "sdk" {
		return &m.SDK, nil
	}
	if key == "version" {
		return &m.Version, nil
	}
	for _, variable := range m.Variables {
		if variable.Key == key {
			return &variable.Value, nil
		}
	}
	return nil, fmt.Errorf("No variable with key '%s' in %+v ", key, m.Variables)
}
