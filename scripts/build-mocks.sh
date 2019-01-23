#!/bin/bash -e

# make sure script is running inside mesg/tools container.
source $(dirname $0)/require-mesg-tools.sh

# generate mocks for container package.
mockery -name=Container -dir ./container -output ./container/mocks

# generate mocks for docker.CommonAPIClient that used by container package.
# TODO: to fix
# mockery -name CommonAPIClient -dir ./vendor/github.com/docker/docker/client -output ./utils/docker/mocks

# generate mocks for database package.
mockery -name=ServiceDB -dir ./database -output ./database/mocks

# generate mocks for commands package.
mockery -name=Executor -dir ./commands -output ./commands/mocks

# generate mocks for protobuf servers.
mockery -name=CoreServer -dir ./protobuf/coreapi -case underscore -inpkg -testonly
mockery -name=ServiceServer -dir ./protobuf/serviceapi -case underscore -inpkg -testonly
