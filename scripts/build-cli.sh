#!/bin/bash -e

if [[ -z "$1" ]]; then
  echo -e "version is not set, run:\n"
  echo "$0 vX.X.X"
  exit 1
fi

LDFLAGS="-s -w -X 'github.com/cosmos/cosmos-sdk/version.Name=mesg' -X 'github.com/cosmos/cosmos-sdk/version.ServerName=mesg-daemon' -X 'github.com/cosmos/cosmos-sdk/version.ClientName=mesg-cli' -X 'github.com/cosmos/cosmos-sdk/version.Version=$1'"
oss=(linux darwin)
archs=(amd64 386)

for os in ${oss[*]}; do
  for arch in ${archs[*]}; do
    echo "Building $os $arch..."
    CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -mod=readonly -o "./bin/mesg-daemon-$os-$arch" -ldflags="$LDFLAGS" ./cmd/mesg-daemon/
    CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -mod=readonly -o "./bin/mesg-cli-$os-$arch" -ldflags="$LDFLAGS" ./cmd/mesg-cli/
  done
done
