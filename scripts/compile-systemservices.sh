#!/bin/bash -e

# compile system services
for s in systemservices/* ; do
  if [ -d "$s" ]; then
    pushd $s > /dev/null
    name=$(basename "$s")
    varname="${name}"
    mesg-cli service:compile | jq -c . > compiled.json
    # LDFLAGS+=" -X 'github.com/mesg-foundation/engine/config.${varname}Compiled=$(cat compiled.json)'"
    popd > /dev/null
  fi
done
