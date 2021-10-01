// +build tools

// Package tools
package tools

import (
	_ "github.com/axw/gocov"
	_ "github.com/mattn/goveralls/tester"
	_ "github.com/modocache/gover/gover"
	_ "github.com/onsi/ginkgo"
	_ "golang.org/x/tools/cover"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
