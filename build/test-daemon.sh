#!/bin/bash

echo "building image mesg/daemon:local..."
docker pull mesg/daemon:latest
docker build -t mesg/daemon:local --cache-from=mesg/daemon:latest --pull .