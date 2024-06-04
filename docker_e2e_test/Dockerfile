# Use the official Ubuntu 20.04 image as the base image
FROM ubuntu:20.04

# Update package list and install necessary dependencies
RUN apt-get update && \
    apt-get dist-upgrade --yes && \
    apt-get install -y build-essential wget zlib1g-dev libncurses5-dev libgdbm-dev libnss3-dev libssl-dev libreadline-dev libffi-dev libsqlite3-dev python3-pip openssh-client


# Download Python 3.6 source code
WORKDIR /tmp
RUN wget https://www.python.org/ftp/python/3.6.9/Python-3.6.9.tgz && \
    tar -xf Python-3.6.9.tgz && \
    cd Python-3.6.9

# Configure, build, and install Python 3.6 from source
RUN cd Python-3.6.9 && \
    ./configure && \
    make && \
    make install

# Verify the installation
RUN python3 --version
COPY . .
RUN python3 -m pip install -r ./requirements.txt
EXPOSE 22 6000 7000
