#!/bin/bash -e

# make sure script is running inside mesg-dev container.
source $(dirname $0)/require-mesg-dev.sh

# navigate to core repo
cd $GOPATH/src/github.com/mesg-foundation/core

# container package
mockery -name=Container -dir ./container -output ./container/mocks

# database package
mockery -name=ServiceDB -dir ./database -output ./database/mocks