package testutils

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

// TestReporting t
func TestReporting(t *testing.T) {
	fakeTest := &testing.T{}
	g := goblin.Goblin(fakeTest)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	reporter := DetailedWithHooksReporter{}
	fakeReporter := goblin.Reporter(&reporter)
	g.SetReporter(fakeReporter)

	g.Describe("One", func() {
		g.It("Foo", func() {
			fmt.Printf("run name it %s", t.Name())
			g.Assert(0).Equal(1)
		})
		g.Describe("Two", func() {
			g.It("Bar", func() {
				g.Assert(0).Equal(0)
			})
		})
	})

	if !reflect.DeepEqual(reporter.Describes, []string{"One", "Two"}) {
		t.FailNow()
	}
	if !reflect.DeepEqual(reporter.Fails, []string{"Foo"}) {
		t.FailNow()
	}
	if !reflect.DeepEqual(reporter.Passes, []string{"Bar"}) {
		t.FailNow()
	}
	if reporter.ends != 2 {
		t.FailNow()
	}

	if !reporter.beginFlag || !reporter.endFlag {
		t.FailNow()
	}
}
