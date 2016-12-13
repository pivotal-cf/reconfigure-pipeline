#!/usr/bin/env bash

set -eux

export GOPATH=$(pwd)/go
export PATH=$GOPATH/bin:$PATH

OUTPUT=$(pwd)/build-binaries-output

VERSION=$(cat version/number)

echo "$VERSION" > "${OUTPUT}/name"
echo "$VERSION" > "${OUTPUT}/tag"

pushd "${GOPATH}/src/github.com/pivotal-cf/reconfigure-pipeline"
  git rev-parse HEAD > "${OUTPUT}/commit"

  go get -v -t ./...

  for os in linux darwin; do
    GOOS=${os} go build

    tar -cvzf "${OUTPUT}/reconfigure-pipeline-${GOOS}.tar.gz" "reconfigure-pipeline"
  done
popd
