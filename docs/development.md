# Development Guide

## Quickstart

The `./build/install-tools.sh` command bootstrap the recommended tools and puts the binaries under ./build/bin
``` shell
./build/install-tools.sh
``` 


## Setup development environment

Following tools are recommended:

| Tool | Usage in the Project | Installation Instructions |
|---|---|---|
| [golangci-lint](https://golangci-lint.run) | linter | [Installation Docs](https://golangci-lint.run/usage/install) </br> `brew install golangci/tap/golangci-lint` </br> `docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint "golangci-lint" run -v`
|
| [go-task](https://taskfile.dev/) | make / build tool | Installation Docs](https://taskfile.dev/#/installation) </br> `brew install go-task/tap/go-task` </br> `sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin` |



