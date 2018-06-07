<p align="center">
  <img src="https://cdn.rawgit.com/mesg-foundation/core/149-update-readme/logo.svg" alt="MESG Core" height="120">
</p>
<h2 align="center">
  <a href="https://mesg.tech/">Website</a> - 
  <a href="https://docs.mesg.tech/">Docs</a> - 
  <a href="https://medium.com/mesg">Blog</a> - 
  <a href="https://discordapp.com/invite/5tVTHJC">Discord</a>
</h2>

<p align="center">
  <a href="https://github.com/mesg-foundation/core"><img src="https://img.shields.io/circleci/project/github/mesg-foundation/core.svg" alt="CircleCI"></a>
  <a href="https://hub.docker.com/r/mesg/core/"><img src="https://img.shields.io/docker/pulls/mesg/core.svg" alt="Docker Pulls"></a>
  <a href="https://codeclimate.com/github/mesg-foundation/core/maintainability"><img src="https://api.codeclimate.com/v1/badges/86ad77f7c13cde40807e/maintainability" alt="Maintainability"></a>
  <a href="https://codecov.io/gh/mesg-foundation/core"><img src="https://codecov.io/gh/mesg-foundation/core/branch/dev/graph/badge.svg" alt="codecov"></a>
</p>

MESG is a platform for the creation of efficient and easy-to-maintain applications that connect any and all technologies.

# Issues

For issues concerning application or service development, please read the [docs](https://docs.mesg.tech/) or ask us directly on [Discord](https://discordapp.com/invite/5tVTHJC) channels #application or #service.

For a question or suggestion of a new feature concerning the Core, please contact us on [Discord](https://discordapp.com/invite/5tVTHJC) channel #core.

To report a bug, please [check for existing issues and create a new issue on this repository](https://github.com/mesg-foundation/core/issues).

# Contribution

For Services and Applications contribution, we have an [curated list of awesome Services and Applications](https://github.com/mesg-foundation/awesome) that you should participate in.

For MESG Core contribution, please contact us on [Discord](https://discordapp.com/invite/5tVTHJC) channel #core. We would love to include you in the development process.

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
env CORE.IMAGE=mesg/core:local go test -cover -v ./...
```

If you use Visual code you can add the following settings (Preference > Settings)
```json
"go.testEnvFile": "${workspaceRoot}/testenv"
```

## Build MESG Core and start it

```bash
./dev-core
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

