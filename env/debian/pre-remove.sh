#!/bin/sh

set -e

MOVIES_SERVICE=movies.service
MOVIES_SOCKET=movies.socket

if [ "$1" = "remove" ]; then
    # stopping service and socket
    if [ -d /run/systemd/system ]; then
        deb-systemd-invoke stop $MOVIES_SERVICE $MOVIES_SOCKET >/dev/null || true
    fi
fi

if [ "$1" = "upgrade" ]; then
    # stopping service
    if [ -d /run/systemd/system ]; then
        deb-systemd-invoke stop $MOVIES_SERVICE >/dev/null || true
    fi
fi
