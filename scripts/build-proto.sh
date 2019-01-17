#!/bin/bash

# make sure script is running inside mesg/tools container.
source $(dirname $0)/require-mesg-tools.sh

PROJECT=/project
GRPC=$PROJECT/protobuf
API_DOCS="--doc_out=$PROJECT/docs/api/ --doc_opt=$PROJECT/docs/api.template"
GRPC_PLUGIN="--go_out=plugins=grpc:./"

protoc $GRPC_PLUGIN $API_DOCS,core.md          --proto_path=$PROJECT $GRPC/coreapi/api.proto
protoc $GRPC_PLUGIN $API_DOCS,service.md       --proto_path=$PROJECT $GRPC/serviceapi/api.proto
