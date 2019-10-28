#!/bin/bash

set -e

# mesg config variables
MESG_NAME=${MESG_NAME:-"engine"}
MESG_PATH=${MESG_PATH:-"$HOME/.mesg"}

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

function start_engine {
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
    --mount type=bind,source=/var/run/docker.sock,destination=/var/run/docker.sock \
    --mount type=bind,source=$MESG_PATH,destination=/root/.mesg \
    --network $MESG_NAME \
    --network name=$MESG_TENDERMINT_NETWORK \
    --publish $MESG_SERVER_PORT:50052 \
    --publish $MESG_TENDERMINT_PORT:26656 \
    mesg/engine:local
}

function stop_engine {
  onexit
}

quiet=false

while getopts "q" o; do
  case $o in
    q)
      quiet=true
      ;;
    *)
      echo "unknown flag $0"
      exit 1
      ;;
  esac
done
shift $((OPTIND-1))

cmd=${1:-"start"}

case $cmd in
  start)
    start_engine
    if ! $quiet; then
      trap onexit EXIT
      docker service logs --follow --raw $MESG_NAME
    fi
    ;;
  stop)
    stop_engine
    ;;
  *)
    echo "unknown command $cmd"
    exit 1
esac
