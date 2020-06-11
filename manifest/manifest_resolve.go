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

func (m Manifest) resolveInstallationInstructions() Instructions {
	return m.Instructions
	/*
	   var result  =     m.Instructions

	   for _, instr := range m.Instructions {
	       switch resolved := instr.(type) {
	       case Pipe:
	           result = append(result, Step{
	               Pipe: instr.ToInstruction(),
	           })
	       case types.DevEnvCommand:
	           result = append(result, Step{resolved})
	       case:
	           result = append(result, Step{Pipe: resolved})
	       }
	   }

	   return result*/
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
