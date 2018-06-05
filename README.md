<center>

![MESG core](https://cdn.rawgit.com/mesg-foundation/core/6be15cc30831d4299d3c3044cb6ad5b5b22c3d86/logo.svg)

# [Website](https://mesg.tech/) - [Docs](https://docs.mesg.tech/) - [Chat](https://discordapp.com/invite/5tVTHJC)

[![CircleCI](https://circleci.com/gh/mesg-foundation/core.svg?style=svg&circle-token=04b7b880e5f42bd26f46a3b11445cb98830e8d92)](https://circleci.com/gh/mesg-foundation/core)
[![Docker Pulls](https://img.shields.io/docker/pulls/mesg/daemon.svg)](https://hub.docker.com/r/mesg/daemon/)
[![Maintainability](https://api.codeclimate.com/v1/badges/86ad77f7c13cde40807e/maintainability)](https://codeclimate.com/github/mesg-foundation/core/maintainability)
[![codecov](https://codecov.io/gh/mesg-foundation/core/branch/dev/graph/badge.svg)](https://codecov.io/gh/mesg-foundation/core)
</center>

## Build from source

### Download source

```bash
mkdir -p $GOPATH/src/github.com/mesg-foundation/core
cd $GOPATH/src/github.com/mesg-foundation/core
git clone https://github.com/mesg-foundation/core.git ./
```

### Install dependencies

```bash
go get -v -t ./...
```

### Run all tests with code coverage

```bash
env DAEMON.IMAGE=mesg/daemon:local go test -cover -v ./...
```

If you use Visual code you can add the following settings (Preference > Settings)
```json
"go.testEnvFile": "${workspaceRoot}/testenv"
```

### Build docker image

```bash
docker build -t mesg-daemon .
```

### Install debugger on OS X

```bash
xcode-select --install
go get -u github.com/derekparker/delve/cmd/dlv
cd $GOPATH/src/github.com/derekparker/delve
make install
```

[Source](https://github.com/derekparker/delve/blob/master/Documentation/installation/osx/install.md)

