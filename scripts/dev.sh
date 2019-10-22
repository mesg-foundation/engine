#!/bin/bash

set -e

# mesg config variables
MESG_NAME=${MESG_NAME:-"engine"}
MESG_PATH=${MESG_PATH:-"$HOME/.mesg"}
MESG_LOG_FORMAT=${MESG_LOG_FORMAT:-"text"}
MESG_LOG_FORCECOLORS=${MESG_LOG_FORCECOLORS:-"false"}
MESG_LOG_LEVEL=${MESG_LOG_LEVEL:-"debug"}
MESG_SERVER_PORT=${MESG_SERVER_PORT:-"50052"}
MESG_TENDERMINT_NETWORK="mesg-tendermint"
MESG_TENDERMINT_PORT=${MESG_TENDERMINT_PORT:-"26656"}

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

mkdir -p $MESG_PATH

echo "create docker service: "
docker service create \
  --name $MESG_NAME \
  --tty \
  --label com.docker.stack.namespace=$MESG_NAME \
  --label com.docker.stack.image=mesg/engine:local \
  --env MESG_NAME=$MESG_NAME \
  --env MESG_LOG_FORMAT=$MESG_LOG_FORMAT \
  --env MESG_LOG_FORCECOLORS=$MESG_LOG_FORCECOLORS \
  --env MESG_LOG_LEVEL=$MESG_LOG_LEVEL \
  --env MESG_TENDERMINT_P2P_PERSISTENTPEERS=$MESG_TENDERMINT_P2P_PERSISTENTPEERS \
  --mount type=bind,source=/var/run/docker.sock,destination=/var/run/docker.sock \
  --mount type=bind,source=$MESG_PATH,destination=/root/.mesg \
  --network $MESG_NAME \
  --network name=$MESG_TENDERMINT_NETWORK \
  --publish $MESG_SERVER_PORT:50052 \
  --publish $MESG_TENDERMINT_PORT:26656 \
  mesg/engine:local

docker service logs --follow --raw $MESG_NAME
