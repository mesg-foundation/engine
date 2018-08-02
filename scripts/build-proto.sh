#!/bin/bash

cd $GOPATH/src

go get github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc

PROJECT=github.com/mesg-foundation/core
CORE=$(pwd)/$PROJECT
API_DOCS="--doc_out=$CORE/docs/api/ --doc_opt=$CORE/docs/api.template"
DATA_DOCS="--doc_out=$CORE/docs/api/ --doc_opt=$CORE/docs/data.template"
GRPC="--go_out=plugins=grpc:./"

protoc $GRPC $DATA_DOCS,service-type.md --proto_path=./ $PROJECT/service/service.proto
protoc $GRPC $API_DOCS,core.md          --proto_path=./ $PROJECT/api/core/api.proto
protoc $GRPC $API_DOCS,service.md       --proto_path=./ $PROJECT/api/service/api.proto
