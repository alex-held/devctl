package plugin_2

import (
	"fmt"
)

func New(args []string) error  {
	return fmt.Errorf("[plugin_2]   arg-len=%d", len(args))
}
