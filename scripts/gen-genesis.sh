#!/bin/bash

# TODO: add some path to be able to run the script from anyfolder
source scripts/dev-config.sh

MESG_VALIDATORS=${MESG_VALIDATORS:-"engine"}

echo "Generating genesis..."
go run internal/tools/gen-genesis/main.go --co-kbpath $MESG_COSMOS_KEYBASE_PATH --tm-path $MESG_TENDERMINT_VALIDATOR_PATH --chain-id $MESG_COSMOS_CHAINID --gentx-filepath $MESG_COSMSOS_GENESISTXFILEPATH --peers-filepath $MESG_COSMSOS_PEERSFILEPATH --validators $MESG_VALIDATORS

echo "Genesis generated with successes and saved to:"
echo " - Cosmos keybase: $MESG_COSMOS_KEYBASE_PATH"
echo " - Tendermint validators: $MESG_TENDERMINT_VALIDATOR_PATH"
echo " - Genesis transaction: $MESG_COSMSOS_GENESISTXFILEPATH"
echo " - Peers: $MESG_COSMSOS_PEERSFILEPATH"
