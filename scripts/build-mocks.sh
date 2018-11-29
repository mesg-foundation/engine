#!/bin/bash -e

# make sure script is running inside mesg/tools container.
source $(dirname $0)/require-mesg-tools.sh

# navigate to core repo.
cd $GOPATH/src/github.com/mesg-foundation/core

# generate mocks for container package.
mockery -name=Container -dir ./container -output ./container/mocks

# generate mocks for docker.CommonAPIClient that used by container package.
mockery -name=Docker -dir ./utils/dockerapi -inpkg

# generate mocks for database package.
mockery -name=ServiceDB -dir ./database -output ./database/mocks

# generate mocks for commands package.
mockery -name=Executor -dir ./commands -output ./commands/mocks
