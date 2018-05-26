#!/bin/bash

docker build -t mesg-daemon-test --cache-from=mesg/daemon --pull .