#!/bin/bash -e

if [[ -z "$1" ]]; then
  echo -e "version is not set, run:\n"
  echo "$0 [version] [prod|dev]"
  exit 1
fi

if [[ -z "$2" ]]; then
  echo -e "release type is not set, can be either prod or dev, run:\n"
  echo "$0 [version] [prod|dev]"
  exit 1
fi


LDFLAGS="-s -w -X 'github.com/mesg-foundation/engine/version.Version=$1'"
archs=(amd64 386)
oss=(linux darwin)

for os in ${oss[*]}; do
  for arch in ${archs[*]}; do
    echo "Building $os $arch..."
    CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -mod=readonly -o "./bin/mesg-daemon-$os-$arch" -ldflags="$LDFLAGS" ./cmd/mesg-daemon/
    CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -mod=readonly -o "./bin/mesg-cli-$os-$arch" -ldflags="$LDFLAGS" ./cmd/mesg-cli/
  done
done

go install github.com/tcnksm/ghr

if [[ "$2" == "dev" ]]; then
  ghr -u mesg-foundation -r engine -delete -prerelease -n "Developer Release" -b "Warning - this is a developer release, use it only if you know what are doing. Make sure to pull the latest \`mesg/engine:dev\` image." release-dev ./bin
fi
if [[ "$2" == "prod" ]]; then
  ghr -u mesg-foundation -r engine -replace $1 ./bin
fi
