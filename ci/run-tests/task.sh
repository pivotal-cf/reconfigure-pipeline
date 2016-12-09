#!/usr/bin/env bash

set -eux

export GOPATH=$(pwd)/go
export PATH=$GOPATH/bin:$PATH

cd $GOPATH/src/github.com/pivotal-cf/reconfigure-pipeline

go get github.com/onsi/ginkgo/ginkgo
go get -v -t ./...

ginkgo -r -race
