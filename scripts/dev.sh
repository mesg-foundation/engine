#!/bin/bash

set -e

MESG_NAME=mesg-engine
MESG_TENDERMINT_NETWORK=mesg-tendermint

function onexit {
  set +e
  echo -e "\nshutting down, please wait..."
  docker_service_remove "$MESG_NAME"
  docker_network_remove "$MESG_NAME"
  docker_network_remove "$MESG_TENDERMINT_NETWORK"
}

function docker_service_remove {
  docker service rm $1
  docker wait $(docker ps -f label=com.docker.swarm.service.name=$1 -q) 2> /dev/null
}

function docker_network_exist {
  [[ ! -z $(docker network list -f name="$1" -f driver=overlay -q) ]]
}

function docker_network_create {
  echo -ne "create docker network $1: "
  docker network create --driver overlay "$1" --label com.docker.stack.namespace="$1"
}

function docker_network_remove {
  echo -ne "remove docker network: "
  docker network remove "$1"
}

trap onexit EXIT

if ! docker_network_exist "$MESG_NAME"; then
  docker_network_create "$MESG_NAME"
fi

if ! docker_network_exist "$MESG_TENDERMINT_NETWORK"; then
  docker_network_create "$MESG_TENDERMINT_NETWORK"
fi

echo "create docker service: "
docker service create \
  --name $MESG_NAME \
  --label com.docker.stack.namespace=$MESG_NAME \
  --label com.docker.stack.image=mesg/engine:dev \
  --mount type=bind,source=/var/run/docker.sock,destination=/var/run/docker.sock \
  --mount type=bind,source="$(git rev-parse --show-toplevel)"/testdata/config/testuser,destination=/root/.mesg/testuser \
  --network $MESG_NAME \
  --network name=$MESG_TENDERMINT_NETWORK \
  --publish 50052:50052 \
  --tty \
  mesg/engine:dev

docker service logs --follow --raw $MESG_NAME
