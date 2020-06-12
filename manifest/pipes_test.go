package manifest

import (
    "testing"

    pipes "github.com/ebuchman/go-shell-pipes"
)

func TestPipe_Execute(t *testing.T) {
    command := "ps aux | grep usr"
    s, err := pipes.RunString(command)
    if err != nil {
        t.Error(err)
    }
    println(s)
}
