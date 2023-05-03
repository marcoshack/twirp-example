//go:build tools
// +build tools

package main

import (
	_ "github.com/securego/gosec/v2/cmd/gosec"
	_ "github.com/twitchtv/twirp/protoc-gen-twirp"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
