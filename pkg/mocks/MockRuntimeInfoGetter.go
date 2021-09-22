package mocks

import "github.com/alex-held/devctl-kit/pkg/system"

type MockRuntimeInfoGetter struct {
	system.RuntimeInfo
}

func (m MockRuntimeInfoGetter) Get() (info system.RuntimeInfo) {
	return m.RuntimeInfo
}
