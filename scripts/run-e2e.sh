#!/bin/bash

set -e

ENGINE_NAME="engine"
NETWORK_NAME="engine"

if [[ -z "$1" ]]; then
  echo -e "arg version is missing, run:\n"
  echo "$0 vX.X.X"
  exit 1
fi


echo "run non existing test to detect compilation error quickly"
go test -mod=readonly -v -count=1 ./e2e/... -run=__NONE__

function onexit {
  docker stop $ENGINE_NAME
  docker network remove $NETWORK_NAME
}

trap onexit EXIT

if [[ -z $(docker network list -f name="$NETWORK_NAME" -q) ]]; then
  docker network create $NETWORK_NAME
fi

docker run \
  -d \
  --rm \
  --name $ENGINE_NAME \
  -p 1317:1317 \
  -p 50052:50052 \
  -p 26657:26657 \
  --network $NETWORK_NAME \
  mesg/engine:$1-dev

echo "waiting lcd server to start"
while true; do
  printf '.'
  block=$(curl --silent http://localhost:1317/node_info | jq .node_info.protocol_version.block)
  if [[ -n $block ]]; then
    break
  fi
  sleep 1
done

echo "starting tests"
go test -failfast -mod=readonly -v -count=1 ./e2e/...
