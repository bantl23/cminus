#!/usr/bin/env bash

go tool yacc -o cminus.y.go -v cminus.y.output cminus.y
go build
./cminus sample0.cm
./cminus sample1.cm
./cminus sample2.cm
./cminus sample3.cm
./cminus sample4.cm
./cminus sample5.cm
./cminus sample6.cm
./cminus sample7.cm
./cminus sample8.cm
./cminus sample9.cm
./cminus sample10.cm
./cminus sample11.cm
