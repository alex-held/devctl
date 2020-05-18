package manifest

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
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
