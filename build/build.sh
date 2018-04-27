#!/usr/bin/env bash

pwd

go get -u github.com/karalabe/xgo

mkdir -p bin
cd bin

xgo \
 -x -v
  --targets=linux/amd64,linux/386,darwin/amd64,windows/amd64,windows/386 \
  ../cli
