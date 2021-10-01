package zsh

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"golang.org/x/sync/syncmap"
)

type resolvable struct {
	ID    string
	Value *string
}

func (r *resolvable) Needs() (needs []string) {
	if r.IsResolved() {
		return needs
	}
	for _, u := range r.Unresolved() {
		needs = append(needs, u)
	}
	return needs
}

func (r *resolvable) NeedsVar(key string) bool {
	if r.IsResolved() {
		return false
	}
	for _, u := range r.Unresolved() {
		if u == key {
			return true
		}
	}
	return false
}

func (c *ResolveContext) LoadEnvironmentVariables() {
	for _, e := range os.Environ() {
		i := strings.Index(e, "=")
		key := e[:i]
		val := e[i+1:]
		fmt.Printf("load %s=%s\n", key, val)
		c.env.Store(key, val)
	}
}

func NewContext(context context.Context) (ctx *ResolveContext) {
	ctx = &ResolveContext{
		env:        syncmap.Map{},
		resolvedC:  make(chan resolvable),
		unresolved: map[string]resolvable{},
		context:    context,
	}

	ctx.cancelC = context.Done()

	// init os envs
	ctx.LoadEnvironmentVariables()

	/*go func() {
		select {
		case r := <-ctx.resolvedC:
			ctx.ResolveWithNeeds(r.ID)
		}
	}()*/
	return ctx
}

// func (r *resolvable) resolve2(ctx *ResolveContext) (fullyResolved bool, resolved *resolvable, unresolved []string) {
// 	if r.IsResolved() {
// 		return true, r, unresolved
// 	}
//
// 	unres := r.Unresolved()
// 	resolvedCount := 0
// 	unresolvedCount := len(unres)
// 	for _, u := range unres {
// 		resolved, ok := ctx.Get(u)
// 		if !ok {
// 			unresolved = append(unresolved, u)
// 			// continue here
// 			unr = append(unr,  ctx.)
// 			continue
// 		}
// 		r.Replace(u, resolved)
// 		resolvedCount++
// 		delete(ctx.unresolved, u)
// 	}
//
// 	for i, i2 := range un {
//
// 	}
// 	ctx.env.Store(r.ID, *r)
// 	return unresolvedCount == resolvedCount, r, unresolved
// }

func (r *resolvable) resolve(ctx *ResolveContext) (ok bool) {
	for !r.IsResolved() {
		unresolved := r.Unresolved()
		maxRetry := len(unresolved)
		for i := 0; i < maxRetry; i++ {
			for _, u := range unresolved {
				resolved, ok := ctx.Get(u)
				if !ok {
					continue
				}
				r.Replace(u, resolved)
			}
		}
		return false
	}

	ctx.env.Store(r.ID, *r.Value)
	ctx.resolvedC <- *r
	return true
}

func (r *resolvable) Name() string {
	return r.ID
}

func (r *resolvable) IsResolved() bool {
	return len(r.Unresolved()) == 0
}

const dollarEscapeToken = "{<{__DOLLAR__}>}"

func escapeDollar(in string) (out string) {
	return strings.ReplaceAll(in, "$$", dollarEscapeToken)
}

func unescapeDollar(in string) (out string) {
	return strings.ReplaceAll(in, dollarEscapeToken, "$$")
}

func (r *resolvable) ReplaceAll(vars map[string]string) {
	for k, v := range vars {
		r.Replace(k, v)
	}
}

func (r *resolvable) Replace(key, val string) {
	reg := getUpdateRegex(key)
	*r.Value = reg.ReplaceAllString(*r.Value, val)
}

const envPattern = `(?:\$)(?:[\{]{0,2})(?P<id>[\w_]*)(?:[\}]{0,2})`

func getUpdateRegex(key string) *regexp.Regexp {
	return regexp.MustCompile(strings.ReplaceAll(envPattern, `(?P<id>[\w_]*)`, fmt.Sprintf("(?P<id>%s)", key)))
}

func (r *resolvable) Unresolved() (unresolved []string) {
	unresolved = []string{}

	reg := regexp.MustCompile(envPattern)

	i := reg.SubexpIndex("id")
	escapedValue := escapeDollar(*r.Value)
	unmatchedVars := reg.FindAllStringSubmatch(escapedValue, -1)
	for _, unmatchedVar := range unmatchedVars {
		unmatched := unmatchedVar[i]
		unescapedMatch := unescapeDollar(unmatched)
		unresolved = append(unresolved, unescapedMatch)
	}

	return unresolved
}

func (r *resolvable) Resolve(ctx *ResolveContext) (ok bool) {
	return r.resolve(ctx)
}
