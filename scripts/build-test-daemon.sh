#!/bin/bash

echo "building image mesg/core:local..."
docker pull mesg/core:latest
docker build -t mesg/core:build --target build .
docker build -t mesg/core:local .