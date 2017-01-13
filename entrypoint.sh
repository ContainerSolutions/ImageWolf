#!/bin/bash
set -euo pipefail
IFS=$'\n\t'
set -e

#remember volumes won't work, which makes things trickier
#could use network pattern I can't remember the name of
#attach to network only works swarm 1.13 and on
docker run -d --name registry-reggie --network reggie -p 5000:5000 \ 
           -e REGISTRY_NOTIFICATIONS_ENDPOINT_URL=${HOSTNAME}:8000/registryNotifications \
           amouat/registry-reggie

exec /reggie $@

