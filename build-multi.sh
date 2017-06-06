#!/bin/bash

set -e

docker build -t imagewolf-build --file build/Dockerfile.build-multi .
docker run -d --name imagewolf-temp imagewolf-build sleep 1h
docker cp imagewolf-temp:/imagewolf-arm ./bin/reggie-arm
docker cp imagewolf-temp:/imagewolf-x86_64 ./bin/reggie-x86_64
docker rm -f imagewolf-temp
docker build -t containersol/imagewolf-armv7l --file build/Dockerfile.run-armv7l .
docker push containersol/imagewolf-armv7l
docker build -t containersol/imagewolf-x86_64 --file build/Dockerfile.run-x86_64 .
docker push containersol/imagewolf-x86_64

manifest pushml ./imagewolf.yaml
