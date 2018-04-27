#!/usr/bin/env bash

go get -u github.com/karalabe/xgo

pwd

cd /go/src/github.com/mesg-foundation/core

mkdir -p bin
cd bin

pwd

echo $GOPATH
GOPATH="/go"
export GOPATH

xgo \
 -x -v \
  --targets=linux/amd64,linux/386,darwin/amd64,windows/amd64,windows/386 \
  /go/src/github.com/mesg-foundation/core/cli
