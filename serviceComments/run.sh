#!/bin/sh

cd `dirname $0`

#while ! nc -z ${DB_HOST} ${DB_PORT}; do sleep 1; done

while true; do
    pg_isready -d multirss -h ${DB_HOST} -p ${DB_PORT} -U postgres && break
    sleep 3
done

test ! -f /lock/comments-schema.lock && \
    psql postgresql://postgres:postgres@${DB_HOST}:${DB_PORT}/multirss < ./db_schema.sql && \
    touch /lock/comments-schema.lock

./serviceComments
