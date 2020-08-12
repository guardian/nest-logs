#!/bin/bash

GOOS=linux GOARCH=amd64 go build main.go
zip nest-logs.zip main