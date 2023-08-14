#!/bin/sh

cd `dirname $0`

while ! nc -z ${DB_HOST} ${DB_PORT}; do sleep 1; done

psql postgresql://postgres:postgres@${DB_HOST}:${DB_PORT}/multirss < ./db_schema.sql
psql postgresql://postgres:postgres@${DB_HOST}:${DB_PORT}/multirss < ./db_data.sql

./serviceNews
