#!/bin/bash

#Run this if issues with exec format
#docker run --rm --privileged multiarch/qemu-user-static:register --reset
set -e

docker build -t imagewolf-build --file build/Dockerfile.build-multi .
docker run -d --name imagewolf-temp imagewolf-build sleep 1h
docker cp imagewolf-temp:/imagewolf-arm ./bin/imagewolf-arm
docker cp imagewolf-temp:/imagewolf-x86_64 ./bin/imagewolf-x86_64
docker rm -f imagewolf-temp
docker build -t containersol/imagewolf-armv7l --file build/Dockerfile.run-armv7l .
docker push containersol/imagewolf-armv7l
docker build -t containersol/imagewolf-x86_64 --file build/Dockerfile.run-x86_64 .
docker push containersol/imagewolf-x86_64

manifest pushml ./imagewolf.yaml
