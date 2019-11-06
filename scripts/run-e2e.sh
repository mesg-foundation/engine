#!/bin/bash

set -e

MESG_PATH="$(pwd)"/e2e.test/mesg

function onexit {
  set +e
  ./scripts/dev.sh stop
  rm -r "${MESG_PATH}"
}

trap onexit EXIT

mkdir -p "${MESG_PATH}"
cp "$(pwd)"/e2e/testdata/e2e.config.yml "${MESG_PATH}"/config.yml

./scripts/dev.sh -q

go test -mod=readonly -v ./e2e/...
