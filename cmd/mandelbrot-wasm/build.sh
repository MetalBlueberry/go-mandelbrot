#!/bin/bash
export GOOS=js
export GOARCH=wasm
go build -o fileserver/lib.wasm main.go