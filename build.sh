#!/bin/bash

set -e

arch=$(uname -m)

docker build -t reggie-build --file Dockerfile.build .
docker run -d --name reggie-temp reggie-build sleep 1h
docker cp reggie-temp:/go/reggie ./reggie-$arch
docker rm -f reggie-temp
docker build --build-arg ARCH=$arch -t amouat/reggie-$arch --file Dockerfile.run .
