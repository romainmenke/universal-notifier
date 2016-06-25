#!/bin/bash

GOOS=linux GOARCH=amd64 go build .

git add .
git commit -m "update binary"
git push
