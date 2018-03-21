[![CircleCI](https://circleci.com/gh/mesg-foundation/application.svg?style=svg&circle-token=04b7b880e5f42bd26f46a3b11445cb98830e8d92)](https://circleci.com/gh/mesg-foundation/application)

# Application

## Build from source

### Download source

```bash
mkdir -p $GOPATH/src/github.com/mesg-foundation/application
cd $GOPATH/src/github.com/mesg-foundation/application
git clone https://github.com/mesg-foundation/application.git ./
```

### Install dependencies

```bash
go get
```

### Run

```bash
go run ./main.go
```

## Install debugger on OS X

```bash
xcode-select --install
go get -u github.com/derekparker/delve/cmd/dlv
cd $GOPATH/src/github.com/derekparker/delve
make install
```

[Source](https://github.com/derekparker/delve/blob/master/Documentation/installation/osx/install.md)

## Run all test with code coverage

```bash
go test -cover ./...
```
