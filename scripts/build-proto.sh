#!/bin/bash

cd $GOPATH/src

PROJECT=github.com/mesg-foundation/core

protoc --go_out=./ $PROJECT/service/service.proto

# build Proto API
protoc --go_out=plugins=grpc:./ --proto_path=./ $PROJECT/api/core/api.proto
protoc --go_out=plugins=grpc:./ --proto_path=./ $PROJECT/api/service/api.proto
