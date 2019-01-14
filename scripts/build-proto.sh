#!/bin/bash

# make sure script is running inside mesg/tools container.
source $(dirname $0)/require-mesg-tools.sh

cd $GOPATH/src

# PROJECT=core
GRPC=/core/protobuf
# CORE=$(pwd)/$PROJECT
API_DOCS="--doc_out=/core/docs/api/ --doc_opt=/core/docs/api.template"
GRPC_PLUGIN="--go_out=plugins=grpc:./"

protoc $GRPC_PLUGIN $API_DOCS,core.md          --proto_path=/core $GRPC/coreapi/api.proto
protoc $GRPC_PLUGIN $API_DOCS,service.md       --proto_path=/core $GRPC/serviceapi/api.proto
