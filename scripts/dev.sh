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
  docker stop $ENGINE_NAME

  if $monitoring; then
    docker stop engine-grafana engine-prometheus
  fi

  docker network remove $NETWORK_NAME
}

trap onexit EXIT

if [[ -z $(docker network list -f name="$NETWORK_NAME" -q) ]]; then
  docker network create $NETWORK_NAME
fi

if $monitoring; then
  echo "start monitoring"
  docker run \
    -d \
    --rm \
    --name=engine-grafana \
    -p 3001:3000 \
    --network $NETWORK_NAME \
    --mount type=bind,source=$(pwd)/scripts/monitoring/datasource.yml,destination=/etc/grafana/provisioning/datasources/datasource.yml \
    --mount type=bind,source=$(pwd)/scripts/monitoring/dashboard.yml,destination=/etc/grafana/provisioning/dashboards/dashboard.yml \
    --mount type=bind,source=$(pwd)/scripts/monitoring/dashboards,destination=/var/lib/grafana/dashboards \
    grafana/grafana

  docker run \
    -d \
    --rm \
    --name=engine-prometheus \
    -p 9090:9090 \
    --network $NETWORK_NAME \
    --mount type=bind,source=$(pwd)/scripts/monitoring/prometheus.yml,destination=/etc/prometheus/prometheus.yml \
    prom/prometheus
fi

docker run \
  -d \
  --rm \
  --name $ENGINE_NAME \
  -p 1317:1317 \
  -p 50052:50052 \
  -p 26657:26657 \
  --network $NETWORK_NAME \
  --volume engine:/root/ \
  mesg/engine:$1-dev

docker logs --tail 1000 --follow $ENGINE_NAME
