#!/bin/sh

find . \( -type d \( -name api -o -name 'vendor*' \) -o -type f -name '*_test.go' \) -prune \
    -o -type f -name "*.go" -exec sh -c 'f={};printf "// === ${f#./} ===\n";cat $f;echo' ';' > merged_go_files.txt
