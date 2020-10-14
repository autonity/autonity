FROM ubuntu:18.04
RUN  apt-get update && \
  apt-get dist-upgrade --yes && \
  apt-get install --yes python3 && \
  apt-get install --yes python3-pip && \
  apt-get install --yes make gcc git wget && \
  wget https://dl.google.com/go/go1.15.2.linux-amd64.tar.gz && \
  tar -xvf go1.15.2.linux-amd64.tar.gz && \
  cp go/bin/* /usr/local/bin && \
  mv go /usr/local && \
  rm go1.15.2.linux-amd64.tar.gz && \
  apt-get install --yes openssh-client
COPY . .

RUN pip3 install -r ./requirements.txt
RUN apt-get clean
EXPOSE 22 6000 7000 8000 9000 30303
