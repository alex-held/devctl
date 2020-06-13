package meta_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"text/template"
	//nolint:gofmt
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/buffer"

	meta2 "github.com/alex-held/dev-env/testdata/meta"
)

func TestRenderTemplate(t *testing.T) {
	meta := meta2.NewDotnetMeta()
	temp := template.New("meta.goyaml")
	temp, err := temp.Parse(meta2.MetaGoyaml)
	if err != nil {
		fmt.Printf("Error while parsing: \n%v", err)
	}
	b := &buffer.Buffer{}
	err = temp.Execute(b, meta)
	if err != nil {
		fmt.Printf("Error while execution: \n%v", err)
		_ = os.Stdout.Sync()
		t.FailNow()
	}
	text := b.String()
	text = strings.Trim(text, "\n")
	_, _ = os.Stdout.WriteString(text)
	assert.NoError(t, err)

}
