#!/bin/bash

source config/env

echo "kill old process ..."
kill -QUIT `cat run/watchmen.pid`

echo "sleep 5 second ..."
sleep 5

echo "start new process ..."
nohup ./watchmen -c config/queues.yml -log log/watchmen.log -pidfile run/watchmen.pid &

exit 0
