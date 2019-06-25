#!/bin/bash

# make sure script is running inside mesg/tools container.
source $(dirname $0)/require-mesg-tools.sh

PROJECT=/project
GRPC=$PROJECT/protobuf
GRPC_PLUGIN="--go_out=plugins=grpc,paths=source_relative:."
protoc $GRPC_PLUGIN --proto_path=$PROJECT $GRPC/definition/execution.proto
protoc $GRPC_PLUGIN --proto_path=$PROJECT $GRPC/definition/service.proto
protoc $GRPC_PLUGIN --proto_path=$PROJECT $GRPC/definition/instance.proto
protoc $GRPC_PLUGIN --proto_path=$PROJECT $GRPC/api/service.proto
protoc $GRPC_PLUGIN --proto_path=$PROJECT $GRPC/api/instance.proto
protoc $GRPC_PLUGIN --proto_path=$PROJECT $GRPC/coreapi/api.proto
protoc $GRPC_PLUGIN --proto_path=$PROJECT $GRPC/api/execution.proto
protoc $GRPC_PLUGIN --proto_path=$PROJECT $GRPC/serviceapi/api.proto
