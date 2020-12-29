FROM ubuntu:18.04
RUN apt-get update && \
    apt-get -y install python3 python3-pip && \
    apt-get clean
RUN pip3 install \
    eth-account~=0.5 \
    eth-utils~=1.9 \
    pytest~=5.4 \
    pytest-xdist~=1.33 \
    requests~=2.25
