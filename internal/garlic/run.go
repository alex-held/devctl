package garlic

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/plugins/plugio"

	"github.com/alex-held/devctl/cli"
)

func isDevctl(mod string) bool {
	if _, err := os.Stat(mod); err != nil {
		return false
	}

	b, err := ioutil.ReadFile(mod)
	if err != nil {
		return false
	}

	return bytes.Contains(b, []byte("github.com/alex-held/devctl"))
}

func Run(ctx context.Context, root string, args []string) error {
	main := filepath.Join(root, "cmd", "devctl")
	if _, err := os.Stat(filepath.Dir(main)); err != nil {
		buff, err := cli.NewFromRoot(root)
		if err != nil {
			return err
		}
		return buff.Main(ctx, root, args)
	}

	bargs := []string{"run", "-v", "./cmd/devctl"}
	bargs = append(bargs, args...)

	cmd := exec.CommandContext(ctx, "go", bargs...)
	cmd.Stdin = plugio.Stdin()
	cmd.Stdout = plugio.Stdout()
	cmd.Stderr = plugio.Stderr()
	return cmd.Run()
}
