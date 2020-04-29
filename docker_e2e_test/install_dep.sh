#!/bin/bash
sudo apt-get update
sudo apt-get install --yes python3
sudo apt-get install --yes python3-pip
sudo apt-get install --yes docker.io
sudo pip3 install -r requirements_docker_test.txt
echo "Dependenices were installed."

