#!/bin/bash

echo "building image mesg/daemon:local..."
docker pull mesg/daemon:latest
docker build -t mesg/daemon:build --target build .
docker build -t mesg/daemon:local .
