#!/bin/bash

# generate relative path of script from core's root.
# scripts can be called from any path including the scripts dir so,
# we need to remove parent paths to get the naked script name.
SCRIPT_PATH=$0
SCRIPT_PATH=${SCRIPT_PATH#*scripts/}
SCRIPT_PATH=./scripts/${SCRIPT_PATH#./}

if [ "$MESG_DEV" != "true" ]; then
   echo "you must run scripts via ./mesg-tools script"
   echo "try executing the following command from the root of core:"
   echo -e "\t./mesg-tools $SCRIPT_PATH $*"
   exit 1
fi
