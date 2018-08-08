#!/bin/bash

echo "Building mesg/core:local..."
docker pull mesg/core:latest
docker build -t mesg/core:build --target build --build-arg version="dev local" .
docker build -t mesg/core:local --build-arg version="dev local" .