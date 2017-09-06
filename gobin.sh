#!/usr/bin/env bash

# install go-bindata
# go get -u github.com/jteeuwen/go-bindata/...
# install goimports
# go get golang.org/x/tools/cmd/goimports

echo "go-bindata resource files..."
go-bindata res/...
goimports -w bindata.go
