#!/bin/sh

set -e

MOVIES_DB_PATH=/var/lib/movies-demo

if [ "$1" = "remove" ]; then
    if [ -f "$MOVIES_DB_PATH/catalog.db" ]; then
        echo "Database file found and won't be removed." >&2
    else
        echo "Removing database folder." >&2
        rm -fr $MOVIES_DB_PATH
    fi

    exit 0
fi
