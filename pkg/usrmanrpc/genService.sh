#!/bin/bash
my_dir=`dirname $0`

cd ../../../../../ &&

protoc \
    -I=$GOPATH/src \
    --gogofaster_out=plugins=grpc:. \
    -I=$GOPATH/src/github.com/gogo/protobuf/protobuf \
    github.com/thanhpp/prom/pkg/usrmanrpc/usrmanService.proto &&

echo "DONE: GEN USRMAN SERVICE PROTO"