#!/bin/sh

set -e

MOVIES_SERVICE=movies.service

# stopping service
if [ -d /run/systemd/system ]; then
    deb-systemd-invoke stop $MOVIES_SERVICE >/dev/null || true
fi
