#!/usr/bin/env bash

go tool yacc -o cminus.y.go -v cminus.y.output cminus.y
go build
./cminus --print-parse-tree --print-symbol-table --print-symbol-map --print-machine-code --trace-analyze --trace-codegen sample0.cm
./cminus --print-parse-tree --print-symbol-table --print-symbol-map --print-machine-code --trace-analyze --trace-codegen sample1.cm
./cminus --print-parse-tree --print-symbol-table --print-symbol-map --print-machine-code --trace-analyze --trace-codegen sample2.cm
./cminus --print-parse-tree --print-symbol-table --print-symbol-map --print-machine-code --trace-analyze --trace-codegen sample3.cm
./cminus --print-parse-tree --print-symbol-table --print-symbol-map --print-machine-code --trace-analyze --trace-codegen sample4.cm
./cminus --print-parse-tree --print-symbol-table --print-symbol-map --print-machine-code --trace-analyze --trace-codegen sample5.cm
./cminus --print-parse-tree --print-symbol-table --print-symbol-map --print-machine-code --trace-analyze --trace-codegen sample6.cm
./cminus --print-parse-tree --print-symbol-table --print-symbol-map --print-machine-code --trace-analyze --trace-codegen sample7.cm
./cminus --print-parse-tree --print-symbol-table --print-symbol-map --print-machine-code --trace-analyze --trace-codegen sample8.cm
./cminus --print-parse-tree --print-symbol-table --print-symbol-map --print-machine-code --trace-analyze --trace-codegen sample9.cm
./cminus --print-parse-tree --print-symbol-table --print-symbol-map --print-machine-code --trace-analyze --trace-codegen sample10.cm
./cminus --print-parse-tree --print-symbol-table --print-symbol-map --print-machine-code --trace-analyze --trace-codegen sample12.cm
./cminus --print-parse-tree --print-symbol-table --print-symbol-map --print-machine-code --trace-analyze --trace-codegen sample13.cm
./cminus --print-parse-tree --print-symbol-table --print-symbol-map --print-machine-code --trace-analyze --trace-codegen sample14.cm
./cminus --print-parse-tree --print-symbol-table --print-symbol-map --print-machine-code --trace-analyze --trace-codegen sample15.cm
./cminus --print-parse-tree --print-symbol-table --print-symbol-map --print-machine-code --trace-analyze --trace-codegen sample16.cm
./cminus --print-parse-tree --print-symbol-table --print-symbol-map --print-machine-code --trace-analyze --trace-codegen sample17.cm
./cminus --print-parse-tree --print-symbol-table --print-symbol-map --print-machine-code --trace-analyze --trace-codegen sample18.cm
./cminus --print-parse-tree --print-symbol-table --print-symbol-map --print-machine-code --trace-analyze --trace-codegen sample19.cm
./cminus --print-parse-tree --print-symbol-table --print-symbol-map --print-machine-code --trace-analyze --trace-codegen sample20.cm
./cminus --print-parse-tree --print-symbol-table --print-symbol-map --print-machine-code --trace-analyze --trace-codegen sample21.cm
./cminus --print-parse-tree --print-symbol-table --print-symbol-map --print-machine-code --trace-analyze --trace-codegen sample22.cm
./cminus --print-parse-tree --print-symbol-table --print-symbol-map --print-machine-code --trace-analyze --trace-codegen sample23.cm
./cminus --print-parse-tree --print-symbol-table --print-symbol-map --print-machine-code --trace-analyze --trace-codegen sample24.cm
./cminus --print-parse-tree --print-symbol-table --print-symbol-map --print-machine-code --trace-analyze --trace-codegen sample25.cm
./cminus --print-parse-tree --print-symbol-table --print-symbol-map --print-machine-code --trace-analyze --trace-codegen sample26.cm
./cminus --print-parse-tree --print-symbol-table --print-symbol-map --print-machine-code --trace-analyze --trace-codegen sample27.cm
./cminus --print-parse-tree --print-symbol-table --print-symbol-map --print-machine-code --trace-analyze --trace-codegen sample28.cm
./cminus --print-parse-tree --print-symbol-table --print-symbol-map --print-machine-code --trace-analyze --trace-codegen sample29.cm
./cminus --print-parse-tree --print-symbol-table --print-symbol-map --print-machine-code --trace-analyze --trace-codegen sample30.cm
