package manifest

import (
	"encoding/json"
	"fmt"
	"github.com/alex-held/dev-env/config"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
	. "path"
	"sort"
	"strings"
)

var DefaultPaths = &config.DefaultPathFactory{}

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

//GetPredefinedVariables resolves global paths. A DefaultPathFactory can be put optionally.
//Only the first factory will be used!
func GetPredefinedVariables(optional ...config.DefaultPathFactory) StringMap {
	factory := config.DefaultPathFactory{}

	if len(optional) > 0 {
		factory = optional[0]
	}
	return StringMap{
		"[[home]]":       factory.GetUserHome(),
		"[[_home]]":      factory.GetDevEnvHome(),
		"[[_sdks]]":      factory.GetSdks(),
		"[[_manifests]]": factory.GetManifests(),
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

	predefinedVariables := GetPredefinedVariables()

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
	predefined := GetPredefinedVariables()
	variables := StringMap{}

	for _, variable := range m.Variables {
		variables[variable.Key] = variable.Value
		key := variable.Key
		value := variable.Value

		// allow manual override of default variables
		if key == "home" || key == "_home" || key == "_sdks" || key == "_installers" || key == "_manifests" {
			predefined[templateStart+key+templateEnd] = value
		}
	}

	for key, value := range variables {

		if strings.HasPrefix(key, templateStart) && strings.HasSuffix(key, templateEnd) {
			continue
		}

		if !strings.HasPrefix(key, templateStart) && !strings.HasSuffix(key, templateEnd) {
			template := templateStart + key + templateEnd
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

//noinspection GoUnusedFunction
func read(fs afero.Fs, path string) (*Manifest, error) {
	manifestRootPath := DefaultPaths.GetManifests()
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
		return nil, fmt.Errorf("Unable to read manifest with path '%s'\nUnknown file extension '%s'\n", manifestPath, fileExtension)
	}
}
