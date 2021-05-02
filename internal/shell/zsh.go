package shell

func Source() string {
	return `
export PATH="$DEVCTL_ROOT/bin:$PATH"
export path=(${(u)path})
	`
}
