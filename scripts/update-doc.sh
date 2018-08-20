#!/bin/bash

echo "Generate API documentation"
./scripts/build-proto.sh
echo "Generate CLI documentation"
go run docs/generate.go