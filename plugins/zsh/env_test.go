package zsh

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolvable_GetUnresolved(t *testing.T) {
	type vars struct {
		value    string
		expected []string
	}
	expectedHomePathKey := "HOME_PATH"

	tts := []vars{
		{
			value:    "$HOME_PATH",
			expected: []string{expectedHomePathKey},
		},
		{
			value:    "${HOME_PATH}",
			expected: []string{expectedHomePathKey},
		},
		{
			value:    "${{HOME_PATH}}",
			expected: []string{expectedHomePathKey},
		},
		{
			value:    "$$HOME_PATH",
			expected: []string{},
		},
		{
			value:    "$HOME_PATH/$HOME_PATH",
			expected: []string{expectedHomePathKey, expectedHomePathKey},
		},
		{
			value:    "$HOME_PATH/$OTHER",
			expected: []string{expectedHomePathKey, "OTHER"},
		},
	}

	ctx := NewContext(context.Background())
	ctx.env.Store("HOME_PATH", "/home/user")

	for _, tt := range tts {
		t.Run("tt.value", func(t *testing.T) {
			sut := resolvable{
				ID:    "foo",
				Value: &tt.value,
			}
			actual := sut.Unresolved()
			t.Logf("%s -> %v", tt.value, actual)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestResolvable_Resolve(t *testing.T) {
	type vars struct {
		value    string
		expected string
	}
	expectedHomePath := "/home/user"

	tts := []vars{
		{
			value:    "$HOME_PATH",
			expected: expectedHomePath,
		},
		{
			value:    "${HOME_PATH}",
			expected: expectedHomePath,
		},
		{
			value:    "${{HOME_PATH}}",
			expected: expectedHomePath,
		},
		{
			value:    "$$HOME_PATH",
			expected: "$$HOME_PATH",
		},
		{
			value:    "$HOME_PATH/$HOME_PATH",
			expected: fmt.Sprintf("%s/%s", expectedHomePath, expectedHomePath),
		},
		{
			value:    "$HOME_PATH/$OTHER",
			expected: fmt.Sprintf("%s/$OTHER", expectedHomePath),
		},
	}

	ctx := NewContext(context.Background())
	ctx.env.Store("HOME_PATH", expectedHomePath)

	for _, tt := range tts {
		t.Run("tt.value", func(t *testing.T) {
			sut := resolvable{
				ID:    "foo",
				Value: &tt.value,
			}
			actual := sut.Resolve(ctx)
			assert.True(t, actual)
			assert.True(t, sut.IsResolved())
			t.Logf("%s -> %v", tt.value, *sut.Value)

			assert.Equal(t, tt.expected, *sut.Value)
		})
	}
}

func TestResolverContext_Resolve(t *testing.T) {
	ctx := NewContext(context.Background())

	valA := "valA"
	valB := "$a $c ${D} $$test"
	valC := "valC"
	valD := "valD"
	valE := "valE"
	valF := "$e $b"

	in := []resolvable{
		{
			ID:    "a",
			Value: &valA,
		},
		{
			ID:    "b",
			Value: &valB,
		},
		{
			ID:    "c",
			Value: &valC,
		},
		{
			ID:    "D",
			Value: &valD,
		},
		{
			ID:    "e",
			Value: &valE,
		},
		{
			ID:    "f",
			Value: &valF,
		},
	}

	ctx.Add(in...)
	ctx.ResolveAll()

	fmt.Println("======== RANGE ========")
	ctx.env.Range(func(key, value interface{}) bool {
		fmt.Printf("%v=%v\n", key, value)
		return true
	})

	actualF, ok := ctx.Get("f")
	assert.True(t, ok)
	assert.Equal(t, "valE valA valC valD $$test", actualF)
}

func TestReplace(t *testing.T) {
	tts := []struct {
		env      map[string]string
		expected string
		name     string
	}{
		{
			name: "with and without brackets",
			env: map[string]string{
				"WITH_BRACKETS": "brackets",
				"NO_BRACKETS":   "no_brackets",
			},
			expected: "$$ESCAPE/brackets/no_brackets",
		},
		{
			name: "with brackets",
			env: map[string]string{
				"WITH_BRACKETS": "brackets",
			},
			expected: "$$ESCAPE/brackets/$NO_BRACKETS",
		},
		{
			name: "without brackets",
			env: map[string]string{
				"NO_BRACKETS": "no_brackets",
			},
			expected: "$$ESCAPE/${WITH_BRACKETS}/no_brackets",
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			v := "$$ESCAPE/${WITH_BRACKETS}/$NO_BRACKETS"
			sut := resolvable{
				ID:    "foo",
				Value: &v,
			}

			for key, vv := range tt.env {
				sut.Replace(key, vv)
			}
			actual := *sut.Value
			assert.Equal(t, tt.expected, actual)
		})
	}
}
