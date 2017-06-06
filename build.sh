#!/bin/bash

set -e

arch=$(uname -m)

docker build -t imagewolf-build --file build/Dockerfile.build .
docker run -d --name imagewolf-temp imagewolf-build sleep 1h
docker cp imagewolf-temp:/go/imagewolf ./bin/imagewolf-$arch
docker rm -f imagewolf-temp
docker build --build-arg ARCH=$arch -t amouat/imagewolf-$arch --file build/Dockerfile.run .
docker push amouat/imagewolf-$arch
