#!/bin/bash

set -e

/registry serve /etc/docker/registry/config.yml &
exec /reggie $@

