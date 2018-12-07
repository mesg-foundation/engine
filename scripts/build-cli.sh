#!/bin/bash -e

if [[ -z "$1" ]]; then
  echo -e "cli version is not set, run:\n"
  echo "$0 cli-version"
  exit 1
fi

LDFLAGS="-X 'github.com/mesg-foundation/core/version.Version=$1'"
archs=(amd64 386)
oss=(linux darwin)

for os in ${oss[*]}; do
  for arch in ${archs[*]}; do
    echo "Building $os $arch..."
    GOOS=$os GOARCH=$arch go build -o "./bin/mesg-core-$os-$arch" -ldflags="$LDFLAGS" ./interface/cli
  done
done
