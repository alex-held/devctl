package manifest

type InstructionType int
type Variables []Variable

type Instructions []Step

type Step struct {
    Command *DevEnvCommand  `json:"command,omitempty"`
    Pipe    []DevEnvCommand `json:"pipe,omitempty"`
}

func (v *Variables) ToMap() StringMap {
	result := StringMap{}

	for _, variable := range *v {
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

type Variable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
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

/*

func (step *Step) ToInstruction() *Pipe {
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
*/
