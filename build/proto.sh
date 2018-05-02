#!/bin/bash

# protoc --go_out=plugins=grpc:./ --proto_path=./types/api types/api/api_service.proto

protoc --go_out=./ service/service.proto