#!/bin/bash -e

LDFLAGS="-X 'github.com/mesg-foundation/engine/version.Version=$version'"

go build -o engine -ldflags="$LDFLAGS" core/main.go
