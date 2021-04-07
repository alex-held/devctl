package shell

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

func (s *Section) RootNode() (root RootCfgNode) { return s.node.RootNode() }

func (s *Section) Execute(w io.Writer, data interface{}) {
	_, _ = fmt.Fprintf(w, "Root Node: %+v; Section Data: %+v", s.RootNode(), data)
}

func NewDevctlSection(data DevCtlSectionConfig) *Section {
	return newSection("devctl", data)
}

func newSection(section string, data interface{}) *Section {
	return &Section{
		node:       nil,
		Title:      strings.ToUpper(section),
		TemplateID: strings.ToLower(section),
		Data:       data,
	}
}

func (s Section) Render() string {
	buf := &bytes.Buffer{}
	cfg := s.getRootConfig()
	t := cfg.Templates.Lookup(s.TemplateID)
	if t == nil {
		panic(fmt.Errorf("failed to lookup required template. template=%s; config=%+v;\n ", s.TemplateID, *cfg))
	}

	err := s.Template.Execute(buf, s.Data)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

type SectionApplyFn func(in *Section) (out *Section)

func (s *Section) Apply(applyFn SectionApplyFn) (out *Section) {
	out = applyFn(s)
	return out
}

type Sections map[string]Section
