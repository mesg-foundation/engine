#!/bin/bash

cd $GOPATH/src

PROJECT=github.com/mesg-foundation/core

protoc --go_out=./ $PROJECT/service/service.proto

for x in $PROJECT/api/*/; do
  protoc \
    --go_out=plugins=grpc:./ \
    --proto_path=./ \
    $x*.proto
done