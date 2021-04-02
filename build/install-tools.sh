#!/usr/bin/env sh
# shellcheck disable=SC2039

set -e

usage() {
	this=$1
	cat << EOF
$this: download tool libraries

Usage: $this [-b] bindir [-d]
  -b sets bindir or installation directory, Defaults to ./build/bin
  -h OR \? displays this help
  -x sets 'set -x' to display all commands executed
EOF
	exit 2
}


parse_args() {
	BINDIR=${BINDIR:-$(pwd)/build/bin}
	while   getopts "b:h?x" arg;do
		case "$arg" in
			b)       BINDIR="$OPTARG";;
			h|        \?) usage "$0";;
			x)       set -x
		esac
	done
	shift   $((OPTIND-1))
}

install_task() {
	if type "$BINDIR"/task &>/dev/null;then
		echo "⏭ tool go-task/task already exists." >&/dev/null
	elif type task &>/dev/null;then
		echo   "⚓️ linking existing go-task/task '$(which task)'"
		ln -s $(which task) $BINDIR/task
	else
		echo   "📦 installing go-task/task https://taskfile.dev/install.sh"
		curl   -sSfL  https://taskfile.dev/install.sh|  sh -s -- -b "$BINDIR"
	fi
}

install_golangcilint() {
	if type "$BINDIR"/golangci-lint &>/dev/null;then
		echo "⏭ tool golangci-lint already exists." >&/dev/null
	elif type golangci-lint &>/dev/null;then
		echo   "⚓️ linking existing golangci-lint '$(which golangci-lint)'"
		ln -s "$(which golangci-lint)" "$BINDIR/golangci-lint"
	else
		echo   "📦 installing golangci-lint https://golangci-lint.run/usage/install"
		curl   -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh|  sh -s -- -b "$BINDIR"
	fi
}

parse_args "$@"


install_task
install_golangcilint
