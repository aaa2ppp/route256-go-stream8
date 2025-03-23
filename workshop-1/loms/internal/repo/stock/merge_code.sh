#!/bin/sh

find . \( -type f -name "*.go" -o -type f -name "*.sql" \) \
    -exec sh -c 'f={};printf "// === ${f#./} ===\n";cat $f;echo' ';'
