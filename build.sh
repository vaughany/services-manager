#!/bin/bash

mkdir -p bin/

echo -e "\e[1mBuilding Linux...\e[0m"
env GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" -a -o bin/ .

echo -e "\e[1mDone.\n\e[0m"
ls -hl bin/

echo
file bin/*