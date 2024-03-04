#!/bin/sh

export GOOS=js
export GOARCH=ecmascript

go list -tags javascript  -f {{.Deps}}
gopherjs build --tags javascript -o md.min.js -m
