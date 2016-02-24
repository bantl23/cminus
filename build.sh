#!/usr/bin/env bash

go get github.com/blynn/nex
nex cminus.nex
go tool yacc cminus.y
go build
