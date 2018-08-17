#!/bin/bash

echo "Building mesg/core:local..."
docker pull mesg/core:latest
docker build -t mesg/core:local --build-arg version="local" .