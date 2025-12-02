#! /bin/bash

TAGS="bindata sqlite sqlite_unlock_notify" make build

if [ -d bin ]; then
    rm -rf bin
fi

mkdir bin
cp kmup bin/
