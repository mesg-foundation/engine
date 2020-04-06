#!/bin/bash

# generate modules' msg
protoc --gogo_out=paths=source_relative:. x/service/internal/types/msg.proto
protoc --gogo_out=paths=source_relative:. x/process/internal/types/msg.proto
protoc --gogo_out=paths=source_relative:. x/runner/internal/types/msg.proto
protoc --gogo_out=paths=source_relative:. x/execution/internal/types/msg.proto
protoc --gogo_out=paths=source_relative:. x/ownership/internal/types/msg.proto

# generate gRPC api
protoc --gogo_out=paths=source_relative,plugins=grpc:. server/grpc/runner/runner.proto
protoc --gogo_out=paths=source_relative,plugins=grpc:. server/grpc/orchestrator/execution.proto
protoc --gogo_out=paths=source_relative,plugins=grpc:. server/grpc/orchestrator/runner.proto
protoc --gogo_out=paths=source_relative,plugins=grpc:. server/grpc/orchestrator/event.proto

TYPES_PATH=protobuf/types
APIS_PATH=protobuf/api

# generate type
for file in "${TYPES_PATH}"/{event,execution,instance,service,process,ownership,runner}.proto
do
  file=$(basename ${file})
  dir="${file%.*}"
  protoc --gogo_out=paths=source_relative:"${dir}" --proto_path . --proto_path "${TYPES_PATH}" "${file}"
done

# generate google type
protoc --gogo_out=paths=source_relative:. protobuf/types/struct.proto

# generate services
for file in "${APIS_PATH}"/{event,execution,instance,service,process,ownership,runner}.proto
do
  protoc --gogo_out=plugins=grpc:. "${file}"
done
