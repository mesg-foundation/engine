#!/bin/bash -e

if [[ -z "$1" ]]; then
  echo -e "version is not set, run:\n"
  echo "$0 [version]"
  exit 1
fi


LDFLAGS="-X 'github.com/mesg-foundation/engine/version.Version=$1'"
archs=(amd64 386)
oss=(linux darwin)

for os in ${oss[*]}; do
  for arch in ${archs[*]}; do
    echo "Building $os $arch..."
    GOOS=$os GOARCH=$arch go build -mod=readonly -o "./bin/mesg-daemon-$os-$arch" -ldflags="$LDFLAGS" ./cmd/mesg-daemon/
    GOOS=$os GOARCH=$arch go build -mod=readonly -o "./bin/mesg-cli-$os-$arch" -ldflags="$LDFLAGS" ./cmd/mesg-cli/
  done
done
