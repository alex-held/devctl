// Code generated by "gonum -types=KindEnum -output=kind_enum.go"; DO NOT EDIT.
// See https://github.com/steinfletcher/gonum
package plugins

import "encoding/json"
import "errors"
import "fmt"

type kindInstanceJsonDescriptionModel struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

var kindInstance = KindEnum{
	SDK: "SDK",
}

// Kind is the enum that instances should be created from
type Kind struct {
	name        string
	value       string
	description string
}

// Enum instances
var SDK = Kind{name: "SDK", value: "SDK", description: "installs updates and manages different sdks on your system"}

// NewKind generates a new Kind from the given display value (name)
func NewKind(value string) (Kind, error) {
	switch value {
	case "SDK":
		return SDK, nil
	default:
		return Kind{}, errors.New(
			fmt.Sprintf("'%s' is not a valid value for type", value))
	}
}

// Name returns the enum display value
func (g Kind) Name() string {
	switch g {
	case SDK:
		return SDK.name
	default:
		return ""
	}
}

// String returns the enum display value and is an alias of Name to implement the Stringer interface
func (g Kind) String() string {
	return g.Name()
}

// Error returns the enum name and implements the Error interface
func (g Kind) Error() string {
	return g.Name()
}

// Description returns the enum description if present. If no description is defined an empty string is returned
func (g Kind) Description() string {
	switch g {
	case SDK:
		return "installs updates and manages different sdks on your system"
	default:
		return ""
	}
}

// KindNames returns the displays values of all enum instances as a slice
func KindNames() []string {
	return []string{
		"SDK",
	}
}

// KindValues returns all enum instances as a slice
func KindValues() []Kind {
	return []Kind{
		SDK,
	}
}

// MarshalJSON provides json serialization support by implementing the Marshaler interface
func (g Kind) MarshalJSON() ([]byte, error) {
	if g.Description() != "" {
		m := kindInstanceJsonDescriptionModel{
			Name:        g.Name(),
			Description: g.Description(),
		}
		return json.Marshal(m)
	}
	return json.Marshal(g.Name())
}

// UnmarshalJSON provides json deserialization support by implementing the Unmarshaler interface
func (g *Kind) UnmarshalJSON(b []byte) error {
	var v interface{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	var value string
	switch v.(type) {
	case map[string]interface{}:
		value = v.(map[string]interface{})["name"].(string)
	case string:
		value = v.(string)
	}

	instance, createErr := NewKind(value)
	if createErr != nil {
		return createErr
	}

	g.name = instance.name
	g.value = instance.value
	g.description = instance.description

	return nil
}