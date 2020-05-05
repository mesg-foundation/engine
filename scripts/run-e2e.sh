#!/bin/bash

set -e

echo "run non existing test to detect compilation error quickly"
go test -mod=readonly -v -count=1 ./e2e/... -run=__NONE__

function onexit {
  docker service rm engine
  docker wait $(docker ps -f label=com.docker.swarm.service.name=engine -q) 2> /dev/null

  docker network remove engine
}

trap onexit EXIT

if [[ -z $(docker network list -f name="engine" -q) ]]; then
  docker network create --driver overlay engine
fi
docker service create \
  --name engine \
  -p 1317:1317 \
  -p 50052:50052 \
  -p 26657:26657 \
  --network engine \
  --label com.docker.stack.namespace=engine \
  mesg/engine:cli-dev

echo "waiting to give some time to the container to start and run"
sleep 10 &
wait $!

echo "starting tests"
go test -failfast -mod=readonly -v -count=1 ./e2e/...
