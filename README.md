<p align="center">
  <img src="https://cdn.rawgit.com/mesg-foundation/core/149-update-readme/logo.svg" alt="MESG Core"
       height="120"><br/>
</p>
<h2 align="center">
  <a href="https://mesg.tech/">Website</a> - 
  <a href="https://docs.mesg.tech/">Docs</a> - 
  <a href="https://medium.com/mesg">Blog</a> - 
  <a href="https://discordapp.com/invite/5tVTHJC">Discord</a>
</h2>

<p align="center">
[![CircleCI](https://img.shields.io/circleci/project/github/mesg-foundation/core.svg)](https://github.com/mesg-foundation/core)
[![Docker Pulls](https://img.shields.io/docker/pulls/mesg/daemon.svg)](https://hub.docker.com/r/mesg/daemon/)
[![Maintainability](https://api.codeclimate.com/v1/badges/86ad77f7c13cde40807e/maintainability)](https://codeclimate.com/github/mesg-foundation/core/maintainability)
[![codecov](https://codecov.io/gh/mesg-foundation/core/branch/dev/graph/badge.svg)](https://codecov.io/gh/mesg-foundation/core)
</p>

**MESG is a platform to help you create efficient and easy to maintain applications that connects any technologie**. You can create your application that listen to an event on a Blockchain and execute a task on a web server. The technology doesn't matter, as long as it can send and/or receive some data.

# Issue

TO DO

# Contribution

TO DO

# Build from source

## Download source

```bash
mkdir -p $GOPATH/src/github.com/mesg-foundation/core
cd $GOPATH/src/github.com/mesg-foundation/core
git clone https://github.com/mesg-foundation/core.git ./
```

## Install dependencies

```bash
go get -v -t -u ./...
```

## Run all tests with code coverage

```bash
env DAEMON.IMAGE=mesg/daemon:local go test -cover -v ./...
```

If you use Visual code you can add the following settings (Preference > Settings)
```json
"go.testEnvFile": "${workspaceRoot}/testenv"
```

## Build daemon and start it

```bash
./dev-daemon
```

## Build CLI and start it

```bash
./dev-cli
```

## Install debugger on OS X

```bash
xcode-select --install
go get -u github.com/derekparker/delve/cmd/dlv
```
If the debugger still doesn't work, try the following:
```bash
cd $GOPATH/src/github.com/derekparker/delve
make install
```

[Source](https://github.com/derekparker/delve/blob/master/Documentation/installation/osx/install.md)

