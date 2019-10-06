#!/bin/bash

# mesg config variables
MESG_NAME=${MESG_NAME:-"engine"}
MESG_PATH=${MESG_PATH:-"$HOME/.mesg"}
MESG_LOG_FORMAT=${MESG_LOG_FORMAT:-"text"}
MESG_LOG_FORCECOLORS=${MESG_LOG_FORCECOLORS:-"false"}
MESG_LOG_LEVEL=${MESG_LOG_LEVEL:-"debug"}
MESG_SERVER_PORT=${MESG_SERVER_PORT:-"50052"}

 # network to connect multiple engine instances
MESG_TENDERMINT_NETWORK="mesg-tendermint"
MESG_TENDERMINT_HOME_CONFIG="$MESG_PATH/tendermint"
MESG_TENDERMINT_VALIDATOR_PORT=${MESG_TENDERMINT_VALIDATOR_PORT:-"26656"}

MESG_COSMOS_GENESISTIME=${MESG_COSMOS_GENESISTIME:-"2019-01-01T00:00:00Z"}
MESG_COSMOS_HOME_CONFIG="$MESG_PATH/cosmos"
MESG_COSMOS_CHAINID=${MESG_COSMOS_CHAINID:-"mesg-chain"}

MESG_TENDERMINT_VALIDATOR_PATH=${MESG_TENDERMINT_VALIDATOR_PATH:-".genesis/tendermint-validator"}
MESG_COSMOS_KEYBASE_PATH=${MESG_COSMOS_KEYBASE_PATH:-".genesis/cosmos-keybase"}
MESG_COSMSOS_GENESISTXFILEPATH=${MESG_COSMSOS_GENESISTXFILEPATH:-".genesis/genesistx.json"}
MESG_COSMSOS_PEERSFILEPATH=${MESG_COSMSOS_PEERSFILEPATH:-".genesis/peers"}
