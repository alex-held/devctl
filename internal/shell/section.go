package shell

import (
	"strings"
)

func (gen *ShellHookGenerator) newSection(section string, data interface{}) *Section {
	return &Section{
		Title:      strings.ToUpper(section),
		TemplateID: strings.ToLower(section),
		GeneratorGetter: func() *ShellHookGenerator {
			return gen
		},
		Data: data,
	}
}

type SectionApplyFn func(in *Section) (out *Section)

func (s *Section) Apply(applyFn SectionApplyFn) (out *Section) {
	out = applyFn(s)
	return out
}

type Sections []Section
