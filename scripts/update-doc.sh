#!/bin/bash

# make sure script is running inside mesg-dev container.
source $(dirname $0)/require-mesg-dev.sh

echo "Generate API documentation"
./scripts/build-proto.sh
echo "Generate CLI documentation"
go run docs/generate.go