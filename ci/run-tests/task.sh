#!/usr/bin/env bash

set -eux

export GOPATH=$(pwd)/go
export PATH=$GOPATH/bin:$PATH

pushd "$GOPATH/src/github.com/pivotal-cf/reconfigure-pipeline"
  go install ./vendor/github.com/onsi/ginkgo/ginkgo
  ginkgo -r -race
popd
