#!/bin/bash

set -e

docker build -t reggie-build --file Dockerfile.build .
docker run -d --name reggie-temp reggie-build sleep 1h
docker cp reggie-temp:/go/reggie ./
docker rm -f reggie-temp
docker build -t reggie --file Dockerfile.run .
