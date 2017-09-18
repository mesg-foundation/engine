#!/bin/bash

set -e

header="Content-Type: application/json"
data='{"jsonrpc":"2.0","method":"net_version","params":[],"id":67}'

until curl -X POST -H "$header" --data "$data" http://parity:8545; do
  >&2 echo "retry in 1s"
  sleep 1
done

>&2 echo "Parity is up - running listener"
exec npm run start
