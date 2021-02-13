package testutils

import (
	"testing"
	
	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestTeardownCombine_Combination_Contains_All_Teardowns(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })
	g.Describe("Teardown", func() {
		
		g.It("WHEN CombineInto is called with one additional teardown => THEN returns new combinator function", func() {
			firstCalled := false
			secondCalled := false
			var first Teardown = func() {
				firstCalled = true
			}
			var second Teardown = func() {
				secondCalled = true
			}
			
			combination := first.CombineInto(second)
			combination()
			
			Expect(firstCalled).To(BeTrue())
			Expect(secondCalled).To(BeTrue())
		})
		
		g.It("WHEN CombineInto is called with multiple additional teardown => THEN returns new combinator function", func() {
			firstCalled := false
			secondCalled := false
			thirdCalled := false
			var first Teardown = func() {
				firstCalled = true
			}
			var second Teardown = func() {
				secondCalled = true
			}
			var third Teardown = func() {
				thirdCalled = true
			}
			
			combination := first.CombineInto(second, third)
			combination()
			
			Expect(firstCalled).To(BeTrue())
			Expect(secondCalled).To(BeTrue())
			Expect(thirdCalled).To(BeTrue())
		})
	})
}
