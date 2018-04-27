#!/usr/bin/env bash

go get -u github.com/karalabe/xgo

mkdir -p bin
cd bin

xgo \
  --targets=linux/amd64,linux/386,darwin/amd64,windows/amd64,windows/386 \
  ../cli
