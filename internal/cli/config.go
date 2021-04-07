package cli

type Config interface {
}

type DefaultConfigImpl struct{}

func NewConfig() Config {
	return &DefaultConfigImpl{}
}
