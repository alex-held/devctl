package shared

import (
	"github.com/alex-held/dev-env/config"
)

func NewTestPathFactory() config.PathFactory {
	homeOverride := "/home"
	return &config.DefaultPathFactory{
		UserHomeOverride: &homeOverride,
		DevEnvDirectory:  ".devenv",
	}
}
