#!/bin/bash

echo Doing mongostuff
/etc/init.d/mongodb start
echo "mongo take its time... therefore sleeping a bit"
sleep 30
mongorestore mongodump/
mongorestore mongodump/

echo "Starting marvin"
supervisord -c supervisord.conf