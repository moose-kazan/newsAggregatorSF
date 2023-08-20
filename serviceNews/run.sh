#!/bin/sh

cd `dirname $0`

while ! nc -z ${DB_HOST} ${DB_PORT}; do sleep 1; done

test ! -f /lock/news-schema.lock && \
    psql postgresql://postgres:postgres@${DB_HOST}:${DB_PORT}/multirss < ./db_schema.sql && \
    touch /lock/news-schema.lock

test ! -f /lock/news-data.lock && \
    psql postgresql://postgres:postgres@${DB_HOST}:${DB_PORT}/multirss < ./db_data.sql && \
    touch /lock/news-data.lock

./serviceNews
