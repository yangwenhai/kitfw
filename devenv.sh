#!/bin/bash
cd `dirname $0`
DIR=`pwd`

GOBIN=$DIR/bin
GOPATH=$DIR


pi=`echo $PATH | grep $GOBIN`
if [ -z "$pi" ]
then
    PATH="$GOBIN:$PATH"
fi

export GOBIN GOPATH PATH
