package plugin_1

import (
	"fmt"
)

func New(args []string) error  {
	return fmt.Errorf("[plugin_1]   arg-len=%d", len(args))
}
