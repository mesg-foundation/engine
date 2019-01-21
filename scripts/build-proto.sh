#!/bin/bash

# make sure script is running inside mesg/tools container.
source $(dirname $0)/require-mesg-tools.sh

PROJECT=/project
GRPC=$PROJECT/protobuf
GRPC_PLUGIN="--go_out=plugins=grpc:./"

protoc $GRPC_PLUGIN --proto_path=$PROJECT $GRPC/coreapi/api.proto
protoc $GRPC_PLUGIN --proto_path=$PROJECT $GRPC/serviceapi/api.proto
