#!/bin/sh

set -e

MOVIES_DB_PATH=/var/lib/movies-demo
MOVIES_USER=movies-demo
MOVIES_SERVICE=movies.service
MOVIES_SOCKET=movies.socket

if [ "$1" = "configure" ]; then
    # creating user and group
    adduser --quiet \
            --system \
            --home /nonexistent \
            --no-create-home \
            --disabled-password \
            --group "$MOVIES_USER"

    # creating database folder
    if [ ! -d $MOVIES_DB_PATH ]; then
        mkdir -p $MOVIES_DB_PATH
        chown $MOVIES_USER:$MOVIES_USER $MOVIES_DB_PATH
    fi

    # enable systemd service
    deb-systemd-helper unmask $MOVIES_SERVICE >/dev/null || true

    if deb-systemd-helper --quiet was-enabled $MOVIES_SERVICE; then
        deb-systemd-helper enable $MOVIES_SERVICE >/dev/null || true
    else
        deb-systemd-helper update-state $MOVIES_SERVICE >/dev/null || true
    fi

    # enable systemd socket
    deb-systemd-helper unmask $MOVIES_SOCKET >/dev/null || true

    if deb-systemd-helper --quiet was-enabled $MOVIES_SOCKET; then
        deb-systemd-helper enable $MOVIES_SOCKET >/dev/null || true
    else
        deb-systemd-helper update-state $MOVIES_SOCKET >/dev/null || true
    fi

    # starting service
    if [ -d /run/systemd/system ]; then
        systemctl --system daemon-reload >/dev/null || true

        if [ -n "$2" ]; then
            deb-systemd-invoke restart $MOVIES_SERVICE >/dev/null || true
        else
            deb-systemd-invoke start $MOVIES_SOCKET >/dev/null || true
        fi

    fi
fi
