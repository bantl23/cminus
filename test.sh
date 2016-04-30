#!/bin/bash

rm *.tm
./cminus --print-parse-tree --optimize=false sample0.cm
mv sample0.tm sample0.none.tm
./cminus --print-parse-tree --optimize=true sample0.cm
mv sample0.tm sample0.opt.tm
./cminus --print-parse-tree --optimize=false sample1.cm
mv sample1.tm sample1.none.tm
./cminus --print-parse-tree --optimize=true sample1.cm
mv sample1.tm sample1.opt.tm
./cminus --print-parse-tree --optimize=false sample2.cm
mv sample2.tm sample2.none.tm
./cminus --print-parse-tree --optimize=true sample2.cm
mv sample2.tm sample2.opt.tm
./cminus --print-parse-tree --optimize=false sample3.cm
mv sample3.tm sample3.none.tm
./cminus --print-parse-tree --optimize=true sample3.cm
mv sample3.tm sample3.opt.tm
./cminus --print-parse-tree --optimize=false sample4.cm
mv sample4.tm sample4.none.tm
./cminus --print-parse-tree --optimize=true sample4.cm
mv sample4.tm sample4.opt.tm
./cminus --print-parse-tree --optimize=false sample5.cm
mv sample5.tm sample5.none.tm
./cminus --print-parse-tree --optimize=true sample5.cm
mv sample5.tm sample5.opt.tm
./cminus --print-parse-tree --optimize=false sample6.cm
mv sample6.tm sample6.none.tm
./cminus --print-parse-tree --optimize=true sample6.cm
mv sample6.tm sample6.opt.tm
./cminus --print-parse-tree --optimize=false sample7.cm
mv sample7.tm sample7.none.tm
./cminus --print-parse-tree --optimize=true sample7.cm
mv sample7.tm sample7.opt.tm
./cminus --print-parse-tree --optimize=false sample8.cm
mv sample8.tm sample8.none.tm
./cminus --print-parse-tree --optimize=true sample8.cm
mv sample8.tm sample8.opt.tm
./cminus --print-parse-tree --optimize=false sample9.cm
mv sample9.tm sample9.none.tm
./cminus --print-parse-tree --optimize=true sample9.cm
mv sample9.tm sample9.opt.tm
./cminus --print-parse-tree --optimize=false sample10.cm
mv sample10.tm sample10.none.tm
./cminus --print-parse-tree --optimize=true sample10.cm
mv sample10.tm sample10.opt.tm
./cminus --print-parse-tree --optimize=false sample11.cm
mv sample11.tm sample11.none.tm
./cminus --print-parse-tree --optimize=true sample11.cm
mv sample11.tm sample11.opt.tm
./cminus --print-parse-tree --optimize=false sample12.cm
mv sample12.tm sample12.none.tm
./cminus --print-parse-tree --optimize=true sample12.cm
mv sample12.tm sample12.opt.tm
