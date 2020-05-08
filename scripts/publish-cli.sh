#!/bin/bash -e

if [[ -z "$1" || -z "$2" || ( "$2" != "unstable" && "$2" != "prod" ) ]]; then
  echo -e "version and release type are not set correctly, run:\n"
  echo "$0 vX.X.X prod|unstable"
  exit 1
fi

go install github.com/tcnksm/ghr

if [[ "$2" == "unstable" ]]; then
  ghr -u mesg-foundation -r engine -p 1 -delete -prerelease -n "Unstable release" -b "Warning - this is an unstable release, use it only if you know what are doing." unstable ./bin
fi
if [[ "$2" == "prod" ]]; then
  ghr -u mesg-foundation -r engine -p 1 -replace "$1" ./bin
fi
