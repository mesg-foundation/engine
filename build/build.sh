#!/usr/bin/env bash

# https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04

package="github.com/mesg-foundation/core/cli"
package_name="mesg-cli"

go get github.com/Microsoft/go-winio
go get github.com/inconshreveable/mousetrap

platforms=("darwin/amd64" "linux/amd64" "linux/386")
# platforms=("windows/amd64" "windows/386" "darwin/amd64" "linux/amd64" "linux/386")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$package_name'-'$GOOS'-'$GOARCH
    
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi  

  mkdir -p ./bin
    env GOOS=$GOOS GOARCH=$GOARCH go build -o ./bin/$output_name $package
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done