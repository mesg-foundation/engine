[![CircleCI](https://circleci.com/gh/mesg-foundation/core.svg?style=svg&circle-token=04b7b880e5f42bd26f46a3b11445cb98830e8d92)](https://circleci.com/gh/mesg-foundation/core)
[![codecov](https://codecov.io/gh/mesg-foundation/core/branch/dev/graph/badge.svg)](https://codecov.io/gh/mesg-foundation/core)

# Core

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
DAEMON.IMAGE=mesg-daemon-test go test -cover -v ./...
```

If you use Visual code you can add the following settings (Preference > Settings)
```json
"go.testEnvFile": "${workspaceRoot}/env.test"
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

