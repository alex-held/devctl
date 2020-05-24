package manifest

import (
	"fmt"
	"sort"
	"strings"
)

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

func (m *Manifest) resolveLinks() []Link {
	variables := m.ResolveVariables()
	var result []Link

	for _, link := range m.Links {
		resolvedSource := ReplaceVariablesIfAny(link.Source, variables)
		resolvedTarget := ReplaceVariablesIfAny(link.Target, variables)
		result = append(result, Link{
			Source: resolvedSource,
			Target: resolvedTarget,
		})
	}

	return result
}

func (m Manifest) resolveInstallationInstructions() []Instructing {
	variables := m.ResolveVariable()
	var result []Instructing

	for _, instr := range m.Instructions {
		instruction := instr.ToInstruction()
		re := instruction.Resolve(variables)
		switch resolved := re.(type) {
		case DevEnvCommand:
			result = append(result, resolved)
		case Pipe:
			result = append(result, resolved)
		}
	}

	return result
}

func (m *Manifest) resolveLinkingInstructions() []Instructing {
	var result []Instructing

	for _, link := range m.resolveLinks() {
		result = append(result, LinkCommand{Link: link})
	}

	return result
}

func (m Manifest) ResolveInstructions() []Instructing {
	var commands []Instructing

	instructions := m.resolveInstallationInstructions()
	linkInstructions := m.resolveLinkingInstructions()

	commands = append(commands, instructions...)
	commands = append(commands, linkInstructions...)

	return commands
}

func (pipe Pipe) Resolve(variables Variables) Instructing {
	return pipe.resolvePipe(variables)
}

func (cmd DevEnvCommand) Resolve(variables Variables) Instructing {
	return cmd.resolveCommand(variables)
}

func (cmd LinkCommand) Resolve(variables Variables) Instructing {
	return cmd.resolveLinkCommand(variables)
}

func (pipe Pipe) resolvePipe(variables Variables) Pipe {
	result := Pipe{Commands: []DevEnvCommand{}}

	for _, command := range pipe.Commands {
		resolvedCommand := command.resolveCommand(variables)
		result.Commands = append(result.Commands, resolvedCommand)
	}
	return result
}

func (cmd LinkCommand) resolveLinkCommand(variables Variables) LinkCommand {
	result := LinkCommand{Link: cmd.Link}

	cmd.Link.Source = ReplaceVariableIfAny(cmd.Link.Source, variables)
	cmd.Link.Target = ReplaceVariableIfAny(cmd.Link.Target, variables)

	return result
}

func (cmd DevEnvCommand) resolveCommand(variables Variables) DevEnvCommand {
	result := DevEnvCommand{
		Command: cmd.Command,
		Args:    []string{},
	}

	for _, commandArg := range cmd.Args {
		resolvedArg := ReplaceVariableIfAny(commandArg, variables)
		result.Args = append(result.Args, resolvedArg)
	}

	return result
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
