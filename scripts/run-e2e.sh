#!/bin/bash

set -e

export MESG_PATH="$(pwd)"/e2e.test/mesg

# first run non existing test to detect compilation error quickly
go test -mod=readonly -v -count=1 ./e2e/... -run=__NONE__

function onexit {
  set +e
  ./scripts/dev.sh stop
  rm -rf "${MESG_PATH}"
}

trap onexit EXIT

mkdir -p "${MESG_PATH}"
cp "$(pwd)"/e2e/testdata/e2e.config.yml "${MESG_PATH}"/config.yml

./scripts/dev.sh -q

go test -mod=readonly -v -count=1 ./e2e/...
