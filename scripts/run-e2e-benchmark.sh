#!/bin/bash

set -e

export MESG_PATH="$(pwd)"/e2e.test/mesg

# first run non existing test to detect compilation error quickly
go test -mod=readonly -v -bench=__NONE__ ./e2e/benchmark...

function onexit {
  set +e
  ./scripts/dev.sh stop
  rm -rf "${MESG_PATH}"
}

trap onexit EXIT

rm -rf "${MESG_PATH}"
mkdir -p "${MESG_PATH}"
cp "$(pwd)"/e2e/testdata/e2e.config.yml "${MESG_PATH}"/config.yml

./scripts/dev.sh -q

go test -failfast -mod=readonly -v -bench=. -benchtime 10x ./e2e/benchmark...
