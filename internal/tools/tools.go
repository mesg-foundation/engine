// +build tools

package tools

// Those imports are here to keep tools in vendor.
// For more info read: https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
import (
	_ "github.com/go-bindata/go-bindata/go-bindata"
	_ "github.com/golang/protobuf/protoc-gen-go"
	_ "github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc"
	_ "github.com/vektra/mockery/cmd/mockery"
)
