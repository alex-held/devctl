package manifest

import (
	"sort"
	"strings"
)

type InstructionType int
type Variables []Variable
type Instructions []Instruction

const (
	Command InstructionType = iota
	CommandPipe
)

type StringSliceStringMap map[string]interface{}
type StringMap map[string]string

type Manifest struct {
	Version      string               `json:"version"`
	SDK          string               `json:"sdk"`
	Variables    StringSliceStringMap `json:"variables"`
	Variable     Variables            `json:"variable"`
	Install      Installer            `json:"install"`
	Instructions Instructions         `json:"instructions"`
	Links        []Link               `json:"links"`
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

func NewInstaller(commands ...Instruction) *Installer {
	i := Installer{
		Instructions: map[int]Instruction{},
		Commands:     []DevEnvCommand{},
	}
	for i2, command := range commands {
		i.Instructions[i2] = command
	}
	return &i
}

type Installer struct {
	Instructions map[int]Instruction `json:"instructions"`
	Commands     []DevEnvCommand     `json:"commands"`
}

type DevEnvCommand struct {
	Command string   `json:"command,omitempty"`
	Args    []string `json:"args,flow"`
}

type Pipe struct {
	Commands []DevEnvCommand `json:"commands,inline"`
}

type Instruction interface {
	Format() string
	Execute() error
	Resolve(variables Variables) Instruction
}

type Variable struct {
	Key   string
	Value string
}
