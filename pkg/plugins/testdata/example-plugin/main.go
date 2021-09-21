package example_plugin

import (
	"bytes"
	"fmt"

	"github.com/alex-held/devctl-plugin/pkg/log"
)

func New(args []string) (err error) {
	fmt.Println("fmt")
	buf := &bytes.Buffer{}
	l := log.New(&log.Config{
		Color: false,
		Out: buf,
		FatalFunc: func() {},
	})


	l.Infof("New called with args=%s", args)
	l.Debugf("Example Plugin name=%s", "New()")
	print(buf.String())

	return fmt.Errorf("some error")
}
