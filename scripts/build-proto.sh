#!/bin/bash

# make sure script is running inside mesg/tools container.
source $(dirname $0)/require-mesg-tools.sh

TYPES_PATH=protobuf/types
APIS_PATH=protobuf/api

# generate type
for file in "${TYPES_PATH}"/{account,event,execution,instance,service,process}.proto
do
  file=$(basename ${file})
  dir="${file%.*}"

  protoc --gogo_out=paths=source_relative:"${dir}" --proto_path . --proto_path "${TYPES_PATH}" "${file}"
done

# generate google type
protoc --gogo_out=paths=source_relative:. protobuf/types/struct.proto

# generate services
for file in "${APIS_PATH}"/{account,event,execution,instance,service,process}.proto
do
  protoc --gogo_out=plugins=grpc:. --proto_path . "${file}"
done
