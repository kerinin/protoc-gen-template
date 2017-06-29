#!/bin/bash

protoc \
    -I=./ \
    -I=/root \
    -I=/root/include \
	--go_out=:$GOPATH/src/ meta/*.proto
