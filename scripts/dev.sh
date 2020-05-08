#!/bin/bash

set -e

ENGINE_NAME="engine"
NETWORK_NAME="engine"

if [[ -z "$1" ]]; then
  echo -e "arg version is missing, run:\n"
  echo "$0 vX.X.X"
  exit 1
fi

monitoring=false
if [[ "$2" == "monitoring" ]]; then
  monitoring=true
fi

function onexit {
  docker service rm $ENGINE_NAME
  docker wait $(docker ps -f label=com.docker.swarm.service.name=$ENGINE_NAME -q) 2> /dev/null

  if $monitor; then
    docker service rm engine-grafana engine-prometheus
    docker wait $(docker ps -f label=com.docker.swarm.service.name=engine-grafana -q) 2> /dev/null
    docker wait $(docker ps -f label=com.docker.swarm.service.name=engine-prometheus -q) 2> /dev/null
  fi

  docker network remove $NETWORK_NAME
}

trap onexit EXIT

if [[ -z $(docker network list -f name="$NETWORK_NAME" -q) ]]; then
  docker network create --driver overlay $NETWORK_NAME
fi

if $monitoring; then
  echo "start monitoring"
  docker service create \
    -p 3001:3000 \
    --network $NETWORK_NAME \
    --name=engine-grafana \
    --mount type=bind,source=$(pwd)/scripts/monitoring/datasource.yml,destination=/etc/grafana/provisioning/datasources/datasource.yml \
    --mount type=bind,source=$(pwd)/scripts/monitoring/dashboard.yml,destination=/etc/grafana/provisioning/dashboards/dashboard.yml \
    --mount type=bind,source=$(pwd)/scripts/monitoring/dashboards,destination=/var/lib/grafana/dashboards \
    grafana/grafana

  docker service create \
    -p 9090:9090 \
    --network $NETWORK_NAME \
    --name=engine-prometheus \
    --mount type=bind,source=$(pwd)/scripts/monitoring/prometheus.yml,destination=/etc/prometheus/prometheus.yml \
    prom/prometheus
fi

docker service create \
  --name $ENGINE_NAME \
  -p 1317:1317 \
  -p 50052:50052 \
  -p 26657:26657 \
  --network $NETWORK_NAME \
  --label com.docker.stack.namespace=$ENGINE_NAME \
  mesg/engine:$1-dev

docker service logs --tail 1000 --follow --raw $ENGINE_NAME
