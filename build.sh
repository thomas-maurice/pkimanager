#!/bin/bash

if ! [ -d bin ]; then mkdir bin; fi;

go get

if ! [ -z "$1" ] && ! [ -z "$2" ]; then
    export GOOS=$1
    export GOARCH=$2
    go build
    exit
fi;

for GOOS in linux; do
    for GOARCH in 386 amd64 arm arm64; do
        echo "Building $GOARCH for system $GOOS"
        export GOOS=$GOOS
        export GOARCH=$GOARCH
        go build -o bin/pkimanager-${GOOS}-$GOARCH
    done
done
