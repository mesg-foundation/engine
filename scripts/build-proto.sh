#!/bin/bash

# make sure script is running inside mesg/tools container.
source $(dirname $0)/require-mesg-tools.sh


PROJECT_PATH=/project
TYPES_PATH=/project/protobuf/types/
APIS_PATH=/project/protobuf/api/
CORE_APIS_PATH=/project/protobuf/coreapi/

# generate types into protobuf/types dir
for t in "${TYPES_PATH}"/{instance,service}.proto
do
  file="$(basename ${t})"
  dir="${file%.*}"
  protoc --gogo_out=Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,plugins=grpc,paths=source_relative:./"${dir}"/ \
    --proto_path=/project/protobuf/types \
    "${file}"
done

# generate types to specific dir
for t in "${TYPES_PATH}"/{event,execution,process}.proto
do
  file="$(basename ${t})"
  dir="${file%.*}"

  protoc --gogo_out=Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,plugins=grpc,paths=source_relative:. \
    --proto_path=${PROJECT_PATH} \
    ${t}
done

# generate services
for t in "${APIS_PATH}"/{event,execution,instance,service,process}.proto
do
  protoc --gogo_out=Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,plugins=grpc,paths=source_relative:. \
    --proto_path=${PROJECT_PATH} \
    "${t}"
done

protoc --go_out=plugins=grpc,paths=source_relative:. \
  --proto_path=${PROJECT_PATH} \
  "${CORE_APIS_PATH}"/api.proto
