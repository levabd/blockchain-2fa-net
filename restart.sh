#!/usr/bin/env bash

# run the net
echo 'restart the net'
docker rm -f $(docker ps -aq) && yes | docker network prune && docker-compose -f network.yaml up -d
