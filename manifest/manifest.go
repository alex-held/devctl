package manifest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alex-held/dev-env/config"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
	. "path"
	"sort"
	"strings"
)

func (m *StringSliceStringMap) ToMap() StringMap {
	var result = map[string]string{}

	for key, val := range *m {
		switch value := val.(type) {
		case string:
			result[key] = value
			continue
		default:
			fmt.Printf("Could not map link of type %t with value %#-v\n", value, value)
			continue
		}
	}

	return result
}

func getPredefinedVariables() StringMap {
	return StringMap{
		"[[home]]":        config.GetUserHome(),
		"[[_home]]":       config.GetDevEnvHome(),
		"[[_sdks]]":       config.GetSdks(),
		"[[_installers]]": config.GetInstallers(),
		"[[_manifests]]":  config.GetManifests(),
	}
}

func getPredefinedVariable() []Variable {
	return []Variable{
		{"[[home]]", config.GetUserHome()},
		{"[[_home]]", config.GetDevEnvHome()},
		{"[[_sdks]]", config.GetSdks()},
		{"[[_installers]]", config.GetInstallers()},
		{"[[_manifests]]", config.GetManifests()},
	}
}

const (
	templateStart = "[["
	templateEnd   = "]]"
)

// ContainsTemplate Returns true if the str contains a [[template]]
func ContainsTemplate(str string) bool { //nolint:whitespace
	startIdx := strings.Index(str, templateStart)
	endIdx := strings.Index(str, templateEnd)
	if startIdx == -1 || endIdx == -1 {
		return false
	}
	return true
}

// Returns the first [[template]] in the str or an error
func GetTemplate(str string) (string, bool) {
	startIdx := strings.Index(str, templateStart)
	endIdx := strings.Index(str, templateEnd) + 2 // add two because of the two characters of ]]
	if ContainsTemplate(str) {
		return str[startIdx:endIdx], true
	}
	return "", false
}

// ResolveTemplateValues adds and resolves all [[template]] in the val
func ResolveTemplateValues(val string, resolved map[string]string) map[string]string { //nolint:whitespace
	if !ContainsTemplate(val) {
		return resolved
	}

	predefinedVariables := getPredefinedVariables()

	for ContainsTemplate(val) {
		template, _ := GetTemplate(val)

		if templateValue, ok := resolved[template]; ok {
			val = strings.ReplaceAll(val, template, templateValue)
			continue
		}

		// try resolve value using predefined variables
		if templateValue, ok := predefinedVariables[template]; ok { // add predefined variable to resolved
			if tValue, ok := resolved[template]; !ok {
				resolved[template] = tValue
			}

			val = strings.ReplaceAll(val, template, templateValue)
			continue
		}

		// try resolve value using previous resolved template values
		if templateValue, ok := resolved[template]; ok {
			if ContainsTemplate(templateValue) {
				resolved = ResolveTemplateValues(val, resolved)
			}

			val = strings.ReplaceAll(val, template, templateValue)
			continue
		}
	}

	return resolved
}

func (m *Manifest) populateVariables() StringMap {
	predefined := getPredefinedVariables()
	variables := StringMap{}

	for _, variable := range m.Variables {
		variables[variable.Key] = variable.Value
	}

	for key, value := range variables {

		if strings.HasPrefix(key, "[[") && strings.HasSuffix(key, "]]") {
			continue
		}

		if !strings.HasPrefix(key, "[[") && !strings.HasSuffix(key, "]]") {
			template := fmt.Sprintf("[[%s]]", key)
			delete(variables, key)
			variables[template] = value
		}
	}

	variables["[[sdk]]"] = m.SDK
	variables["[[version]]"] = m.Version

	for key, value := range predefined {
		variables[key] = value
	}

	return variables
}

func (sm *StringMap) SortByKeys() []string {
	keys := make([]string, 0, len(*sm))

	for k := range *sm {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func ReplaceVariableIfAny(s string, variables Variables) string {
	for _, variable := range variables {
		if ContainsTemplate(s) {
			s = strings.ReplaceAll(s, variable.Key, variable.Value)
		}
	}
	return s
}
func ReplaceVariablesIfAny(s string, variables map[string]string) string {
	for key, value := range variables {
		if ContainsTemplate(s) {
			s = strings.ReplaceAll(s, key, value)
		}
	}
	return s
}

func (m *Manifest) ResolveVariables() StringMap {
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

	return variables
}

func (m *Manifest) VariableMap() map[string]string {
	return m.Variables.ToMap()
}

/*
func (m Manifest) ResolveCommands() []Instructing {
    variables := m.ResolveVariables()
    var result []Instructing

    for _, command := range m.Instructions {
        switch c := command.(type) {
        case DevEnvCommand:

            devEnvCommand := DevEnvCommand{
                Command: c.Command,
                Args:    []string{},
            }

            for _, arg := range c.Args {
                resolvedArguments := ReplaceVariablesIfAny(arg, variables)
                devEnvCommand.Args = append(devEnvCommand.Args, resolvedArguments)
            }

            result = append(result, &devEnvCommand)

        case Pipe:

        }

    }

    return result
}
*/
/*func (i Step) Format() string {
    switch i.Type {
    case Command:
        command := DevEnvCommand{
            Command: i.Command,
            Args:    i.Args,
        }
        return command.Format()
    case CommandPipe:
        pipe := Pipe{
            Commands: i.Commands,
        }
        return pipe.Format()
    default:
        return fmt.Sprintf("%+v", i)
    }
}*/

func (pipe Pipe) Format() string {
	sb := strings.Builder{}
	maxIdx := len(pipe.Commands) - 1

	for idx, command := range pipe.Commands {
		formatted := command.Format()
		sb.WriteString(formatted)
		if idx < maxIdx {
			sb.WriteString(" | ")
		}
	}

	return sb.String()
}

func (cmd LinkCommand) GetCommands() []DevEnvCommand {

	return []DevEnvCommand{
		{
			Command: "mkdir",
			Args:    []string{"-p", Dir(cmd.Link.Target)},
		},
		{
			Command: "ln",
			Args:    []string{"-s", cmd.Link.Source, cmd.Link.Target},
		},
	}

}

func (cmd LinkCommand) Format() string {
	sb := strings.Builder{}
	commands := cmd.GetCommands()
	for _, command := range commands {
		sb.WriteString(fmt.Sprintf("%s; ", command.Format()))
	}
	return sb.String()
}

func (cmd DevEnvCommand) Format() string {
	command := cmd.Command
	for _, arg := range cmd.Args {
		command += fmt.Sprintf(" %s", arg)
	}
	return command
}

/*
func (i Step) Execute() error {
    switch i.Type {
    case Command:
        command := DevEnvCommand{
            Command: i.Command,
            Args:    i.Args,
        }
        return command.Execute()
    case CommandPipe:
        pipe := Pipe{
            Commands: i.Commands,
        }
        return pipe.Execute()
    default:
        return fmt.Errorf("Invalid instruction type: '%v' ", i.Type)
    }
}*/

func readJson(text string, manifest *Manifest) (*Manifest, error) {
	err := json.Unmarshal([]byte(text), manifest)

	if err != nil {
		return manifest, err
	}

	return manifest, nil
}

func readYaml(text string, manifest *Manifest) (*Manifest, error) {
	err := yaml.Unmarshal([]byte(text), manifest)

	if err != nil {
		return manifest, err
	}

	return manifest, nil
}

func read(fs afero.Fs, path string) (*Manifest, error) {
	manifestRootPath := config.GetManifests()
	manifestPath := Join(manifestRootPath, path)
	file, err := afero.ReadFile(fs, manifestPath)
	fileExtension := Ext(manifestPath)
	manifest := &Manifest{}

	if err != nil {
		return nil, err
	}

	switch fileExtension {
	case ".json":
		return readJson(string(file), manifest)
	case ".yaml":
		return readYaml(string(file), manifest)
	default:
		return nil, errors.New(fmt.Sprintf("Unable to read manifest with path '%s'\nUnknown file extension '%s'\n", manifestPath, fileExtension))
	}
}
