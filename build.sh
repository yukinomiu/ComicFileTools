#!/bin/bash

mainDir=`cd src/command && pwd`
outDir=`cd bin && pwd`

function build() {
  CGO_ENABLED=0 GOOS=${1} GOARCH=${2} go build -trimpath -ldflags "-s -w -buildid=" -o ${3}/comic-file-tools-${1}-${2}
}

cd $mainDir
build darwin amd64 $outDir
build windows amd64 $outDir

