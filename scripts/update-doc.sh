#!/bin/bash

# make sure script is running inside mesg-tools container.
source $(dirname $0)/require-mesg-tools.sh

echo "Generate API documentation"
./scripts/build-proto.sh

echo "Generate CLI documentation"
go run docs/generate.go
