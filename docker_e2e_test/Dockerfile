# Use the official Ubuntu 24.04 image as the base image for docker test engine.
FROM ubuntu:24.04

# Update package list and install necessary dependencies
RUN apt-get update && \
    apt-get dist-upgrade --yes && \
    apt-get install -y build-essential wget zlib1g-dev libncurses5-dev libgdbm-dev libnss3-dev libssl-dev libreadline-dev libffi-dev libsqlite3-dev python3-pip openssh-client

# Verify the installation
RUN python3 --version
RUN python3 -m pip install -r ./requirements.txt
EXPOSE 22 6000 7000
