#!/bin/bash

set -e

# mesg config variables
MESG_NAME=${MESG_NAME:-"bob"}
MESG_PATH=${MESG_PATH:-"$HOME/.mesg"}
MESG_LOG_FORMAT=${MESG_LOG_FORMAT:-"text"}
MESG_LOG_FORCECOLORS=${MESG_LOG_FORCECOLORS:-"false"}
MESG_LOG_LEVEL=${MESG_LOG_LEVEL:-"debug"}
MESG_SERVER_PORT=${MESG_SERVER_PORT:-"50052"}

MESG_VALIDATOR_NUMBER=${MESG_VALIDATOR_NUMBER:-"1"}

 # network to connect multiple engine instances
MESG_TENDERMINT_NETWORK="mesg-tendermint"
MESG_TENDERMINT_HOME_CONFIG="$MESG_PATH/tendermint"
MESG_TENDERMINT_VALIDATOR_PORT=${MESG_TENDERMINT_VALIDATOR_PORT:-"26656"} # to delete or should be passed to engine config to start multiple port at the same time maybe
MESG_TENDERMINT_VALIDATOR_PATH=${MESG_TENDERMINT_VALIDATOR_PATH:-".genesis/tendermint-validator"}

MESG_COSMOS_HOME_CONFIG="$MESG_PATH/cosmos"
MESG_COSMOS_CHAINID=${MESG_COSMOS_CHAINID:-"mesg-chain"}
MESG_COSMOS_KEYBASE_PATH=${MESG_COSMOS_KEYBASE_PATH:-".genesis/cosmos-keybase"}
MESG_COSMOS_GENESISTIME=${MESG_COSMOS_GENESISTIME:-"2019-01-01T00:00:00Z"}
MESG_COSMSOS_GENESISTXFILEPATH=${MESG_COSMSOS_GENESISTXFILEPATH:-".genesis/genesistx.json"}
MESG_COSMSOS_PEERSFILEPATH=${MESG_COSMSOS_PEERSFILEPATH:-".genesis/peers"}

if [[ $* == *--gen-genesis* ]]; then
  echo "Flag --gen-genesis passed. Generating new genesis..."
  go run internal/tools/gen-genesis/main.go --co-kbpath $MESG_COSMOS_KEYBASE_PATH --tm-path $MESG_TENDERMINT_VALIDATOR_PATH --chain-id $MESG_COSMOS_CHAINID --gentx-filepath $MESG_COSMSOS_GENESISTXFILEPATH --peers-filepath $MESG_COSMSOS_PEERSFILEPATH --vno $MESG_VALIDATOR_NUMBER
fi

MESG_COSMOS_GENESISVALIDATORTX=$(cat $MESG_COSMSOS_GENESISTXFILEPATH)
MESG_TENDERMINT_P2P_PERSISTENTPEERS=$(cat $MESG_COSMSOS_PEERSFILEPATH)

# Setup the validator private keys
if [[ $* == *--validator* ]]; then
  echo "Flag --validator passed. Copy validator private keys to $MESG_TENDERMINT_HOME_CONFIG"
  mkdir -p $MESG_TENDERMINT_HOME_CONFIG
  rsync -a $MESG_TENDERMINT_VALIDATOR_PATH/$MESG_NAME/ $MESG_TENDERMINT_HOME_CONFIG
  MESG_TENDERMINT_VALIDATOR_PORT_PUBLISH="--publish $MESG_TENDERMINT_VALIDATOR_PORT:26656"
fi

# TODO: to remove when account importation is implemented
if [[ $* == *--genesis_account ]]; then
  echo "Flag --genesis_account passed. Copy genesis account private keys to $MESG_COSMOS_HOME_CONFIG"
  mkdir -p $MESG_COSMOS_HOME_CONFIG
  rsync -a $MESG_COSMOS_KEYBASE_PATH/ $MESG_COSMOS_HOME_CONFIG
fi

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
  --tty \
  --label com.docker.stack.namespace=$MESG_NAME \
  --label com.docker.stack.image=mesg/engine:dev \
  --env MESG_NAME=$MESG_NAME \
  --env MESG_LOG_FORMAT=$MESG_LOG_FORMAT \
  --env MESG_LOG_FORCECOLORS=$MESG_LOG_FORCECOLORS \
  --env MESG_LOG_LEVEL=$MESG_LOG_LEVEL \
  --env MESG_TENDERMINT_P2P_PERSISTENTPEERS=$MESG_TENDERMINT_P2P_PERSISTENTPEERS \
  --env MESG_TENDERMINT_P2P_EXTERNALADDRESS=$MESG_TENDERMINT_P2P_EXTERNALADDRESS \
  --env MESG_COSMOS_CHAINID=$MESG_COSMOS_CHAINID \
  --env MESG_COSMOS_GENESISVALIDATORTX=$MESG_COSMOS_GENESISVALIDATORTX \
  --env MESG_COSMOS_GENESISTIME=$MESG_COSMOS_GENESISTIME \
  --mount type=bind,source=/var/run/docker.sock,destination=/var/run/docker.sock \
  --mount type=bind,source=$MESG_PATH,destination=/root/.mesg \
  --network $MESG_NAME \
  --network name=$MESG_TENDERMINT_NETWORK \
  --publish $MESG_SERVER_PORT:50052 \
  $MESG_TENDERMINT_VALIDATOR_PORT_PUBLISH \
  mesg/engine:dev

docker service logs --follow --raw $MESG_NAME