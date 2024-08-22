#!/bin/zsh
gofmt -s -w .
export $(cat .env | xargs)
go build -o ./cmd src/main.go
./cmd/main
