#!/bin/bash -e

# go get mockery
go get github.com/vektra/mockery/...

# navigate to core repo
cd $GOPATH/src/github.com/mesg-foundation/core

# container package
mockery -name=Container -dir ./container -output ./container/mocks