#!/usr/bin/env bash

export SAWTOOTH_CORE=/home/peshkov/dev/blockchain/sawtooth-core
docker-compose -f network_two_validators.yaml up
