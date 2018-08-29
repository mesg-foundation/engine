#!/bin/bash -e

go get ./...
mkdir -p bin

archs=(amd64 386)
oss=(linux darwin)

for os in ${oss[*]}; do
  for arch in ${archs[*]}; do
    echo "Building $os $arch..."
    GOOS=$os GOARCH=$arch go build \
      -o ./bin/mesg-core-$os-$arch \
      -ldflags="-X 'github.com/mesg-foundation/core/version.Version=$1'" \
      ./interface/cli
  done
done
