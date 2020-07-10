package shared

func NewTestPathFactory() PathFactory {
	homeOverride := "/Users/dev/.go/src/github.com/alex-held/dev-env/testdata"
	return &DefaultPathFactory{
		UserHomeOverride: &homeOverride,
		DevEnvDirectory:  ".devenv",
	}
}
