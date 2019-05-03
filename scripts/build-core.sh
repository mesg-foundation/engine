#!/bin/bash -e

MESG_SERVICE_SERVER=ipfs.app.mesg.com

LDFLAGS="-X 'github.com/mesg-foundation/core/version.Version=$version'"
LDFLAGS+=" -X 'github.com/mesg-foundation/core/config.EnvMarketplaceEndpoint=https://mainnet.infura.io/v3/7690bc6d35e140d2be4e771a1237f636'"
LDFLAGS+=" -X 'github.com/mesg-foundation/core/config.EnvMarketplaceAddress=0x0C6e8d0eC4770fDa8A56CD912392d2ff14822952'"
LDFLAGS+=" -X 'github.com/mesg-foundation/core/config.EnvMarketplaceToken=0x420167d87d35c3a249b32ef6225872fbd9ab85d2'"

# upload system services
for s in systemservices/* ; do
  if [ -d "$s" ]; then
    pushd $s > /dev/null
    name=$(basename "$s")
    varname="${name^}"
    archive="$name.tar.gz"

    id=$(
      tar -czf - --exclude-from=.dockerignore . |
      curl -s -F "file=@-;filename=${archive}" http://$MESG_SERVICE_SERVER:5001/api/v0/add |
      jq -r .Hash
    )
    LDFLAGS+=" -X 'github.com/mesg-foundation/core/config.${varname}URL=http://$MESG_SERVICE_SERVER:8080/ipfs/$id'"
    popd > /dev/null
  fi
done

go build -o mesg-core -ldflags="$LDFLAGS" core/main.go
