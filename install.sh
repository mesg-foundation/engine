#!/bin/bash
# Copyright 2018 mesg.tech

ARCH=$(uname -m)
if [[ "$OSTYPE" == "linux-gnu" && "$ARCH" == "x86_64" ]]; then
  # Linux amd64
  BINARY_LINK=https://github.com/mesg-foundation/core/releases/download/release-dev/mesg-core-linux-amd64
  BINARY_LOCAL=/usr/local/bin/mesg-core
elif [[ "$OSTYPE" == "linux-gnu" ]]; then
  # Linux i386
  BINARY_LINK=https://github.com/mesg-foundation/core/releases/download/release-dev/mesg-core-linux-386
  BINARY_LOCAL=/usr/local/bin/mesg-core
elif [[ "$OSTYPE" == "darwin"* && "$ARCH" == "x86_64" ]]; then
  # Mac OSX x86_64
  BINARY_LINK=https://github.com/mesg-foundation/core/releases/download/release-dev/mesg-core-darwin-10.6-amd64
  BINARY_LOCAL=~/.local/bin/mesg-core
elif [[ "$OSTYPE" == "darwin"* ]]; then
  # Mac OSX i386
  BINARY_LINK=https://github.com/mesg-foundation/core/releases/download/release-dev/mesg-core-darwin-10.6-386
  BINARY_LOCAL=~/.local/bin/mesg-core
# elif [[ "$OSTYPE" == "cygwin" ]]; then
#   # POSIX compatibility layer and Linux environment emulation for Windows
# elif [[ "$OSTYPE" == "msys" ]]; then
#   # Lightweight shell and GNU utilities compiled for Windows (part of MinGW)
# elif [[ "$OSTYPE" == "win32" ]]; then
#   # I'm not sure this can happen.
# elif [[ "$OSTYPE" == "freebsd"* ]]; then
#   # ...
else
  # Unknown.
  echo "Automatic installation is not compatible with your system yet. Please download the binary from https://github.com/mesg-foundation/core/releases"
  exit
fi

echo "Download binary from GitHub"
curl $BINARY_LINK --progress-bar -L -o $BINARY_LOCAL
chmod +x $BINARY_LOCAL

echo "Installation finished."
echo "To run the MESG Core, execute: mesg-core"
