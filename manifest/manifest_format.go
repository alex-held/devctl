package manifest

import (
	"encoding/json"
	"fmt"
	"github.com/disiqueira/gotree"
	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"
	"strings"
)

type tableFormat struct {
	Caption string
	Table   *tablewriter.Table
	Writer  *strings.Builder
}

type FormatType int

const (
	Table FormatType = iota
	Tree
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

func (m *Manifest) Format(formatType FormatType) string {
	switch formatType {
	case Table:
		return m.FormatTable()
	case Tree:
		return m.FormatAsTree()
	default:
		return fmt.Sprintf("%+v", *m)
	}
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

	for _, link := range m.resolveLinks() {
		links.Add(fmt.Sprintf("%+v", link))
	}
	/*
	   for idx, cliExe := range m.resolveInstallationInstructions() {
	       instruction := gotree.New(fmt.Sprintf("%d", idx))
	       switch command := cliExe.(type) {
	       case types.DevEnvCommand:
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
	*/
	root.AddTree(variable)
	root.AddTree(links)
	root.AddTree(instructions)
	formattedTree := root.Print()
	return formattedTree
}

func newTableFormat(name string, appender func(appender *tablewriter.Table), headers ...string) tableFormat {
	writer := &strings.Builder{}
	return tableFormat{
		Caption: name,
		Table:   createTable(writer, appender, headers...),
		Writer:  writer,
	}
}

func (m *Manifest) formatTables() []tableFormat {
	properties := newTableFormat("Properties", func(table *tablewriter.Table) {
		table.Append([]string{"Version", m.Version})
		table.Append([]string{"SDK", m.SDK})
	}, "Property", "Value")

	variables := newTableFormat("Variables", func(table *tablewriter.Table) {
		for _, variable := range m.ResolveVariable() {
			table.Append([]string{variable.Key, variable.Value})
		}
	}, "Variables", "Value")

	links := newTableFormat("Links", func(table *tablewriter.Table) {
		for _, link := range m.resolveLinks() {
			table.Append([]string{link.Source, link.Target})
		}
	}, "Source", "Target")

	instructions := newTableFormat("Instructions", func(table *tablewriter.Table) {
		for _, step := range m.GetInstructions() {
			table.Append([]string{fmt.Sprintf("%d", step)})
		}
	}, "Order", "Command")

	return []tableFormat{
		properties,
		variables,
		links,
		instructions,
	}
}

func (tableFormat *tableFormat) Format() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("\n"))
	sb.WriteString(fmt.Sprintf("%-3s\n", tableFormat.Caption))
	tableString := tableFormat.Writer.String()
	sb.WriteString(tableString)
	return sb.String()
}

func createTable(writer *strings.Builder, appender func(appender *tablewriter.Table), headers ...string) *tablewriter.Table {
	table := tablewriter.NewWriter(writer)
	table.SetHeader(headers)
	/*
	   table.SetHeaderColor(
	       tablewriter.Colors{tablewriter.Bold, tablewriter.FgWhiteColor},
	       tablewriter.Colors{tablewriter.Bold, tablewriter.FgGreenColor},
	   )

	   table.SetColumnColor(
	       tablewriter.Colors{tablewriter.FgWhiteColor},
	       tablewriter.Colors{tablewriter.FgGreenColor},
	   )
	*/
	table.SetAutoFormatHeaders(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderLine(true)
	table.SetAutoWrapText(false)
	table.SetBorder(false)
	table.SetCenterSeparator(" ")
	table.SetColumnSeparator("|")
	table.SetRowSeparator("-")
	table.SetAutoMergeCells(true)

	appender(table)
	table.Render()
	return table
}

// Formats the Source into a colorful table representation
func (m *Manifest) FormatTable() string {
	sb := strings.Builder{}
	for _, table := range m.formatTables() {
		sb.WriteString(table.Format())
	}
	return sb.String()
}
