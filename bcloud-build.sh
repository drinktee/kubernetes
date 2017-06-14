#!/bin/bash
set -eu
set -o pipefail

# log control
: ${SCIRPT:="$0"}

MODULE=kubernetes-bce

# source path
WORKROOT=$(pwd)

export GOROOT=/home/scmtools/buildkit/go/go_1.7.3/
export GO=$GOROOT/bin/go
export PATH=$PATH:$GOROOT/bin


function log() {
    echo "$(date +%F_%T) [${SCIRPT}][Notice] : $*"
}

function log2exit() {
    echo "$(date +%F_%T) [${SCIRPT}][Error] : $*"
    exit 1
}

function build()
{
    log "building..."
    make
}

build $*

### del unneeded files
# find -name ".git" -prune -exec rm -rf {} \;
# find -name ".svn" -prune -exec rm -rf {} \;

### tar all files
mv $WORKROOT/_output $WORKROOT/output
# tar zcf data.tar.gz ./
# mkdir -p $GOPATH/output/
# mv data.tar.gz $GOPATH/output/