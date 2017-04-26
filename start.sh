#!/bin/bash

echo Doing mongostuff
/etc/init.d/mongodb start
echo "mongo takes its time... therefore sleeping a bit"
sleep 30
#mongorestore mongodump/

mongoimport --db quotes --collection quote --file quotes.json --jsonArray

echo "Starting marvin"
supervisord -c supervisord.conf