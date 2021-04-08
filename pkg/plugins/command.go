package plugins

type Command struct {
	// Name "foo"
	Name string `json:"name"`

	// UseCommand "bar"
	UseCommand string `json:"use_command"`

	// DevCtlCommand "sdk"
	DevCtlCommand string `json:"devctl_command"`

	// Description "manages the foo sdk"
	Description string   `json:"description,omitempty"`
	Aliases     []string `json:"aliases,omitempty"`
	Binary      string   `json:"-"`
	Flags       []string `json:"flags,omitempty"`
	// Filters events to listen to ("" or "*") is all events
	ListenFor string `json:"listen_for,omitempty"`
}

type Commands []Command
