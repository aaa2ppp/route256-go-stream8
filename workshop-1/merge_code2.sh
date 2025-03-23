#!/bin/sh

# find . \( -type d \( -name api -o -name 'vendor*' \) -o -type f -name '*_test.go' \) -prune \
#     -o -type f -name "*.go" -exec sh -c 'f={};printf "// === ${f#./} ===\n";cat $f;echo' ';' > merged_go_files.txt

find . \( -type d -name 'vendor*' \) -prune \
    -o \( \
        -type f -name 'Dockerfile' \
        -o -type f -name 'docker-compose.yml' \
        -o -type f -name 'Makefile' \
        -o -type f -name '*.mk' \
        -o -type f -name '*.sh' \
        -o -type f -name '*.proto' \
    \) \
    -exec sh -c 'f={};printf "// === ${f#./} ===\n";cat $f;echo' ';'