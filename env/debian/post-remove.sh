#!/bin/sh

set -e

MOVIES_DB_PATH=/var/lib/movies-demo
MOVIES_SERVICE=movies.service

if [ "$1" = "remove" ]; then
    if [ -f "$MOVIES_DB_PATH/catalog.db" ]; then
        echo "Database file found and won't be removed." >&2
    else
        echo "Removing database folder." >&2
        rm -fr $MOVIES_DB_PATH
    fi

    # disabling service
    if [ -d /run/systemd/system ]; then
        systemctl --system daemon-reload >/dev/null || true
    fi

    deb-systemd-helper mask $MOVIES_SERVICE >/dev/null || true
fi

if [ "$1" = "purge" ]; then
    # disabling service
    deb-systemd-helper purge $MOVIES_SERVICE >/dev/null || true
    deb-systemd-helper unmask $MOVIES_SERVICE >/dev/null || true
fi
