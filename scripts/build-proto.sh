#!/bin/bash

# make sure script is running inside mesg/tools container.
source $(dirname $0)/require-mesg-tools.sh

cd $GOPATH/src

PROJECT=github.com/mesg-foundation/core
GRPC=$PROJECT/protobuf
CORE=$(pwd)/$PROJECT
API_DOCS="--doc_out=$CORE/docs/api/ --doc_opt=$CORE/docs/api.template"
GRPC_PLUGIN="--go_out=plugins=grpc:./"

protoc $GRPC_PLUGIN $API_DOCS,core.md          --proto_path=./ $GRPC/coreapi/api.proto
protoc $GRPC_PLUGIN $API_DOCS,service.md       --proto_path=./ $GRPC/serviceapi/api.proto
