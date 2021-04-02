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

log_info() {
	tool=$1
	printf '[install-tools]\tInstalling %s 📦\n\n' "$tool"
}

install_task() {
	if  (type task &>/dev/null)|| (type "$BINDIR"/task &>/dev/null);then
		log_info   "📦 installing go-task/task https://taskfile.dev/install.sh"
		curl   -sSfL  https://taskfile.dev/install.sh|  sh -s -- -b "$BINDIR"
	fi
}

install_golangcilint() {
	if  (type golangci-lint &>/dev/null)|| (type "$BINDIR"/golangci-lint &>/dev/null);then
		log_info   "📦 installing go-task/task https://golangci-lint.run/usage/install"
		curl   -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh|  sh -s -- -b "$BINDIR"
	fi
}

parse_args "$@"

install_task
install_golangcilint
