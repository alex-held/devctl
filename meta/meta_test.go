package meta_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"text/template"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/buffer"

	"github.com/alex-held/dev-env/shared"
	meta2 "github.com/alex-held/dev-env/testdata/meta"
)

func init() {
	shared.BootstrapLogger(zerolog.TraceLevel)
}
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
	fmt.Println(text)
	assert.NoError(t, err)
}

func TestRenderPackage(t *testing.T) {
	meta := meta2.NewDotnetMeta()
	meta.Sources = nil
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
	log.Trace().Str("output", text).Msg("Rendered Packages")
	assert.NoError(t, err)
}
