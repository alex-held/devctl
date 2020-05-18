package manifest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alex-held/dev-env/cmd/install"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
	path2 "path"
	"sort"
	"strings"
)

type StringSliceStringMap map[string]interface{}
type StringMap map[string]string

type Manifest struct {
	Version   string               `json:"version"`
	SDK       string               `json:"sdk"`
	Variables StringSliceStringMap `json:"variables"`
	Install   Install              `json:"install"`
	Links     []Link               `json:"links"`
}

type Link struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

// Formats the Manifest into a colorful table representation
func (m *Manifest) Format() string {
	variables := m.ResolveVariables()
	tableString := &strings.Builder{}

	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"Kind", "Key", "Value"})

	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Normal, tablewriter.BgHiBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.BgRedColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.BgHiGreenColor},
	)

	table.SetColumnColor(
		tablewriter.Colors{tablewriter.Normal, tablewriter.BgWhiteColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgRedColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgGreenColor},
	)

	table.SetAutoFormatHeaders(true)

	table.SetHeaderLine(true)
	table.SetAutoWrapText(true)
	table.SetBorder(false)
	table.SetCenterSeparator("+")
	table.SetColumnSeparator("|")
	table.SetRowSeparator("-")

	table.SetTablePadding("\t")
	table.SetAutoMergeCells(true)

	var vKeys []string
	for key := range variables {
		vKeys = append(vKeys, key)
	}
	sort.Strings(vKeys)

	for _, key := range vKeys {
		table.Append([]string{"Variable", key, variables[key]})
	}
	table.Append([]string{"", "", ""})

	for _, val := range m.Links {

		table.Append([]string{"Link", val.Source, val.Target})
	}

	table.Render()
	return tableString.String()
}

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
		"[[home]]":       install.GetUserHome(),
		"[[_home]]":      install.GetDevEnvHome(),
		"[[_sdks]]":      install.GetSdks(),
		"[[_installers]": install.GetInstallers(),
		"[[_manifests]]": install.GetManifests(),
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
	variables := m.Variables.ToMap()

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

type Install struct {
	Commands []string `json:"commands"`
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

func read(fs afero.Fs, path string) (*Manifest, error) {
	manifestRootPath := install.GetManifests()
	manifestPath := path2.Join(manifestRootPath, path)
	file, err := afero.ReadFile(fs, manifestPath)
	fileExtension := path2.Ext(manifestPath)
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
