#!/bin/bash

# script is use to use local go cache to speed up build
# docker doesn't allow to mount volume during build, 
# so it always rebuild whole project. Use go build cache to 
# make docker build faster.

set -e

ENGINE_SUM_PATH="./bin/.engine.sum"
DOCKER_SUM_PATH="./bin/.Dockerfile.dev.sum"

echo "compile engine for linux amd64 with CGO disabled"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/engine-linux-amd64 -ldflags="-s -w -X 'github.com/mesg-foundation/engine/version.Version=local'" core/main.go

touch "$ENGINE_SUM_PATH" "$DOCKER_SUM_PATH"

# check if engine bin was cached
ENGINE_SUM="$(openssl md5 ./bin/engine-linux-amd64)"
ENGINE_SUM_PREV="$(cat $ENGINE_SUM_PATH)"
if [[ "$ENGINE_SUM" == "$ENGINE_SUM_PREV" ]]; then
  BINCACHED=true
else
  echo "$ENGINE_SUM" > "$ENGINE_SUM_PATH"
fi

# check if dockerfile was cached
DOCKER_SUM="$(openssl md5 ./Dockerfile.dev)"
DOCKER_SUM_PREV="$(cat $DOCKER_SUM_PATH)"
if [[ "$DOCKER_SUM" == "$DOCKER_SUM_PREV" ]]; then
  DOCKERCACHED=true
else
  echo "$DOCKER_SUM" > "$DOCKER_SUM_PATH"
fi

IMAGE_EXIST="$(docker image ls mesg/engine:local -q)"
if [[ ! $BINCACHED ]] || [[ ! $DOCKERCACHED ]] || [[ $IMAGE_EXIST == "" ]]; then
  echo "build mesg/engine image"
  docker build -f Dockerfile.dev -t "mesg/engine:local" .
fi
