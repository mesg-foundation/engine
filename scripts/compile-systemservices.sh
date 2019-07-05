#!/bin/bash -e

# compile system services
for s in systemservices/* ; do
  if [ -d "$s" ]; then
    pushd $s > /dev/null
    echo "Compiling $s"
    mesg-cli service:compile | jq -c . > compiled.json
    popd > /dev/null
  fi
done
