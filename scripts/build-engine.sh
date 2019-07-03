#!/bin/bash -e

LDFLAGS="-X 'github.com/mesg-foundation/engine/version.Version=$version'"
LDFLAGS+=" -X 'github.com/mesg-foundation/engine/config.EnvMarketplaceEndpoint=https://mainnet.infura.io/v3/7690bc6d35e140d2be4e771a1237f636'"
LDFLAGS+=" -X 'github.com/mesg-foundation/engine/config.EnvMarketplaceAddress=0x0C6e8d0eC4770fDa8A56CD912392d2ff14822952'"
LDFLAGS+=" -X 'github.com/mesg-foundation/engine/config.EnvMarketplaceToken=0x420167d87d35c3a249b32ef6225872fbd9ab85d2'"

for s in systemservices/* ; do
  if [ -d "$s" ]; then
    pushd $s > /dev/null
    name=$(basename "$s")
    varname="${name^}"
    LDFLAGS+=" -X 'github.com/mesg-foundation/engine/config.${varname}Compiled=$(cat compiled.json)'"
    popd > /dev/null
  fi
done

go build -o engine -ldflags="$LDFLAGS" core/main.go
