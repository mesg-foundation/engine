#!/bin/bash -e

# make sure script is running inside mesg/tools container.
source $(dirname $0)/require-mesg-tools.sh

# generate mocks
mockery -name Container -dir ./container -output ./container/mocks
mockery -name CommonAPIClient -dir ./internal/mocks -output ./utils/docker/mocks
