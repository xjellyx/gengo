#!/bin/bash
dir=$GOPATH
echo $dir

go build -o gengo ./cmd/gengo/main.go

mv gengo $dir/bin

echo "install success"