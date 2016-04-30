#!/bin/bash

./cminus --print-parse-tree --optimize=false sample12.cm
mv sample12.tm sample12.tm-none
./cminus --print-parse-tree --optimize=true sample12.cm
mv sample12.tm sample12.tm-opt
