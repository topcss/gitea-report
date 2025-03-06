#!/bin/bash

# 清理
rm -rf release

VERSION="1.0.1"
APPNAME="gitea-report"

# 编译 windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o release/$APPNAME-$VERSION-windows-amd64.exe

# 编译 linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o release/$APPNAME-$VERSION-linux-amd64

# 编译 arm64
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o release/$APPNAME-$VERSION-linux-arm64
