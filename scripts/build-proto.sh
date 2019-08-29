#!/bin/bash

# make sure script is running inside mesg/tools container.
source $(dirname $0)/require-mesg-tools.sh

PROJECT=/project
GRPC=$PROJECT/protobuf
GRPC_PLUGIN="--go_out=plugins=grpc,paths=source_relative:."

protoc --gogofaster_out=plugins=grpc,paths=source_relative:./event/ \
       --proto_path=/project/protobuf/types \
       event.proto

protoc --gogofaster_out=plugins=grpc,paths=source_relative:./execution/ \
       --proto_path=/project/protobuf/types \
       execution.proto

protoc --gogofaster_out=plugins=grpc,paths=source_relative:./instance/ \
       --proto_path=/project/protobuf/types \
       instance.proto

protoc --gogofaster_out=plugins=grpc,paths=source_relative:./workflow/ \
       --proto_path=/project/protobuf/types \
       workflow.proto

protoc --gogofaster_out=plugins=grpc,paths=source_relative:./service/ \
       --proto_path=/project/protobuf/types \
       service.proto

protoc $GRPC_PLUGIN --proto_path=$PROJECT $GRPC/api/event.proto
protoc $GRPC_PLUGIN --proto_path=$PROJECT $GRPC/api/execution.proto
protoc $GRPC_PLUGIN --proto_path=$PROJECT $GRPC/api/instance.proto
protoc $GRPC_PLUGIN --proto_path=$PROJECT $GRPC/api/service.proto
protoc $GRPC_PLUGIN --proto_path=$PROJECT $GRPC/api/workflow.proto

protoc $GRPC_PLUGIN --proto_path=$PROJECT $GRPC/coreapi/api.proto
