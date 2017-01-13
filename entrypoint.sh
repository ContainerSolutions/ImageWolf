#!/bin/bash
set -euo pipefail
IFS=$'\n\t'
set -e

sigterm()
{
   echo "Shutting down"
   docker stop registry-reggie
   exit 0
}

trap 'sigterm' TERM

#remember volumes won't work, which makes things trickier
#could use network pattern I can't remember the name of
#attach to network only works swarm 1.13 and on
docker pull amouat/registry-reggie
#kill any old registry
docker stop registry-reggie || true
docker rm registry-reggie || true
docker run -d --name registry-reggie --network reggie -p 5000:5000 \
           -e REGISTRY_NOTIFICATIONS_ENDPOINTS_URL=${HOSTNAME}:8000/registryNotifications \
           amouat/registry-reggie

/reggie $@ &

#note that this is started with tini, so shouldn't need to pass signal to reggie
while /usr/bin/true ; do
  sleep 30
done