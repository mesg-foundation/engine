#!/bin/bash

go get ./...
go get github.com/karalabe/xgo
mkdir bin
cd bin; xgo --targets=darwin/386,darwin/amd64,linux/386,linux/amd64,windows/386,windows/amd64 -ldflags="-X main.version=$1" -out mesg-core ../cli
sudo chmod +x ./bin/*
