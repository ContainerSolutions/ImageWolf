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
docker rm -f registry-reggie || true # -f needed for cases where can't be stopped e.g fs error
export REGISTRY_NOTIFICATIONS_ENDPOINTS=$(cat <<EOF
    - name: reggie
      disabled: false
      url: http://${HOSTNAME}:8000/registryNotifications
EOF
)
docker run -d --name registry-reggie --network reggie -p 5000:5000 \
           -e REGISTRY_NOTIFICATIONS_ENDPOINTS \
           amouat/registry-reggie

/reggie $@ &

#note that this is started with tini, so shouldn't need to pass signal to reggie
while true ; do
  sleep 30
done
