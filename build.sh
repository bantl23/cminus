#!/usr/bin/env bash

rm *.tm-*
rm cminus
go tool yacc -o cminus.y.go -v cminus.y.output cminus.y
go build
