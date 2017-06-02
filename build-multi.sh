#!/bin/bash

set -e

docker build -t reggie-build --file build/Dockerfile.build-multi .
docker run -d --name reggie-temp reggie-build sleep 1h
docker cp reggie-temp:/reggie-arm ./bin/reggie-arm
docker cp reggie-temp:/reggie-x86_64 ./bin/reggie-x86_64
docker rm -f reggie-temp
docker build -t amouat/reggie-armv7l --file build/Dockerfile.run-armv7l .
docker push amouat/reggie-armv7l
docker build -t amouat/reggie-x86_64 --file build/Dockerfile.run-x86_64 .
docker push amouat/reggie-x86_64

manifest pushml ./reggie.yaml
