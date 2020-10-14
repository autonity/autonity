FROM ubuntu:18.04
RUN  apt-get update && \
  apt-get dist-upgrade --yes && \
  apt-get install --yes python3 && \
  apt-get install --yes python3-pip && \
  #apt-get install --yes golang make gcc musl-dev libc-dev git && \
  apt-get install --yes golang make gcc git wget && \
  apt-get install --yes openssh-client
COPY . .
RUN git clone https://github.com/clearmatics/autonity.git && cd ./autonity && git checkout master && make all
COPY ./autonity/build/bin/autonity ./bin/
COPY ./autonity/build/bin/bootnode ./bin/
RUN rm -rf ./autonity

RUN pip3 install -r ./requirements.txt
RUN apt-get clean
EXPOSE 22 6000 7000 8000 9000 30303
