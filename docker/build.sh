#!/bin/bash
if [[ $1 == 'amd' ]]; then
    docker buildx build --platform linux/amd64 .
else
    docker buildx build .
fi
