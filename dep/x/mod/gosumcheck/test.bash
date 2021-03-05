#!/bin/bash

set -e
go build -o gosumcheck.exe
export GONOSUMDB=*/text # rsc.io/text but not utilware/dep/x/text
./gosumcheck.exe "$@" -v test.sum
rm -f ./gosumcheck.exe
echo PASS
