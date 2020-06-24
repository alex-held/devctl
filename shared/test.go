package shared

func NewTestPathFactory() PathFactory {
	homeOverride := "/home"
	return &DefaultPathFactory{
		UserHomeOverride: &homeOverride,
		DevEnvDirectory:  ".devenv",
	}
}
