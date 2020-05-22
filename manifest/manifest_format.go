package manifest

import (
	"encoding/json"
	"fmt"
	"github.com/disiqueira/gotree"
	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"
	"sort"
	"strings"
)

func PrintYaml(manifest interface{}) {
	var y, _ = yaml.Marshal(manifest)
	fmt.Println("YAML")
	fmt.Println(string(y))
}

func PrintJson(manifest interface{}) {
	j, _ := json.MarshalIndent(manifest, "", "  ")
	fmt.Println("JSON")
	fmt.Println(string(j))
}

func (m *Manifest) FormatAsTree() string {
	root := gotree.New(fmt.Sprintf("%s-%s", m.SDK, m.Version))

	variable := gotree.New("variables")
	links := gotree.New("links")
	instructions := gotree.New("instructions")

	for _, value := range m.ResolveVariable() {
		formatted := fmt.Sprintf("%+v", value)
		variable.Add(formatted)
	}

	for _, link := range m.ResolveLinks() {
		links.Add(fmt.Sprintf("%+v", link))
	}

	for idx, cliExec := range m.ResolveInstructions() {
		instruction := gotree.New(fmt.Sprintf("%d", idx))

		switch command := cliExec.(type) {
		case DevEnvCommand:
			formatted := fmt.Sprintf("%s", command.Format())
			instruction.Add(formatted)
		case Pipe:
			for _, command := range command.Commands {
				formatted := fmt.Sprintf("%s", command.Format())
				instruction.Add(formatted)
			}
		}

		instructions.AddTree(instruction)
	}

	root.AddTree(variable)
	root.AddTree(links)
	root.AddTree(instructions)
	formattedTree := root.Print()
	return formattedTree
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
