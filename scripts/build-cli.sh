#!/bin/bash

go get ./...
go get github.com/karalabe/xgo
mkdir bin
xgo --targets=darwin/386,darwin/amd64,linux/386,linux/amd64,windows/386,windows/amd64 \
  -ldflags="-X 'github.com/mesg-foundation/core/version.Version=$1'" \
  -out mesg-core \
  -dest ./bin \
  ./cli
sudo chmod +x ./bin/*
