#!/bin/sh

cd `dirname $0`

while ! nc -z ${DB_HOST} ${DB_PORT}; do sleep 1; done

./serviceNews
