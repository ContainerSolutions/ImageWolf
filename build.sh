#!/bin/bash

set -e

# TODO: check uname -m and use to build arm/x86_64 image

docker build -t reggie-build --file Dockerfile.build .
docker run -d --name reggie-temp reggie-build sleep 1h
docker cp reggie-temp:/go/reggie ./
docker rm -f reggie-temp
docker build -t reggie --file Dockerfile.run .
