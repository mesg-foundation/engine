#!/usr/bin/env bash

package="github.com/mesg-foundation/core/cli"
package_name="mesg-cli"

go get -u github.com/karalabe/xgo

mkdir -p bin
cd bin

xgo \
  --deps=https://gmplib.org/download/gmp/gmp-6.1.0.tar.bz2 \
  --targets=linux/amd64,linux/386,darwin/amd64,windows/amd64,windows/386 \
  ../cli
