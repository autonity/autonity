# Test engine docker file
FROM ubuntu:22.04
RUN  apt-get update && \
  apt-get dist-upgrade --yes && \
  apt-get install --yes python3 && \
  apt-get install --yes python3-pip && \
  apt-get install --yes openssh-client
COPY . .
RUN pip3 install -r ./requirements.txt
RUN apt-get clean
EXPOSE 22 6000 7000 8000 30303 20202
