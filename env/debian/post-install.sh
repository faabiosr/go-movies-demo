#!/bin/sh

set -e

MOVIES_DB_PATH=/var/lib/movies-demo
MOVIES_USER=movies-demo

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

    exit 0
fi
