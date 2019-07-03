#!/usr/bin/env bash

(cd ~/go/src/github.com/clearmatics/autonity/ && make all)

python3.6 main.py ~/go/src/github.com/clearmatics/autonity/build/bin/autonity