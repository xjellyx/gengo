#!/bin/bash
dir=$GOPATH
echo $dir

go build -o gengo ./cmd/gengo.go

mv gengo $dir/bin

echo "install success"