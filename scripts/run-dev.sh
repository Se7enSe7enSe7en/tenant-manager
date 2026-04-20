#!/usr/bin/env bash

# Build the development binary
go build -tags=debug -o dist/main-dev ./cmd/server

# If build succeeds, neatly replace this bash process with the newly compiled binary 
# so it directly receives all stop signals from templ
if [ $? -eq 0 ]; then
    exec ./dist/main-dev
fi
