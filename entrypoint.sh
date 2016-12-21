#!/bin/bash

set -e

/registry serve /etc/docker/registry/config.yml &
/reggie $@ &
sleep infinity

