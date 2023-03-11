#!/bin/bash

# macOS
GOOS=darwin GOARCH=amd64 go build -o bin/app-amd64-mac main.go
GOOS=darwin GOARCH=arm64 go build -o bin/app-arm64-mac main.go

# windows
GOOS=windows GOARCH=amd64 go build -o bin/app-amd64-windows.exe main.go

# linux
GOOS=linux GOARCH=amd64 go build -o bin/app-amd64-linux main.go