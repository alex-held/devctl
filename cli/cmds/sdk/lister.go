package sdk

type ListerCmd struct{}

func (c *ListerCmd) PluginName() string {
	return "sdk/lister"
}
