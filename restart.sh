#!/usr/bin/env bash

# build projects
echo 'build allatrack/tfa-entrance-backend image'
cd ../blockchain-2fa-backend
docker build -t allatrack/tfa-entrance-backend .

echo 'build allatrack/tfa-entrance-frontend image'
cd ../blockchain-2fa-frontend
docker build -t allatrack/tfa-entrance-frontend .

echo 'build allatrack/tfa-cabinet image'
cd ../blockchain-2fa-web
docker build -t allatrack/tfa-cabinet .

cd ../blockchain-2fa-net
echo 'build allatrack/sawtooth-tfa-s-tp image'
cd ./service_tfa_processor
docker build -t allatrack/sawtooth-tfa-s-tp .

echo 'build allatrack/sawtooth-tfa-sc-tp image'
cd ../client_tfa_processor
docker build -t allatrack/sawtooth-tfa-sc-tp .

# run the net
echo 'run the net'
cd ..
docker rm -f $(docker ps -aq) && yes | docker network prune && docker-compose -f network.yaml up -d
