# Test engine docker file.
FROM jrei/systemd-ubuntu
RUN  apt-get update && \
  apt-get dist-upgrade --yes && \
  apt-get install --yes python3 && \
  apt-get install --yes python3-pip && \
  apt-get install --yes openssh-client
RUN mkdir /var/run/sshd
RUN echo 'root:root' |chpasswd
RUN mkdir /root/.ssh
RUN useradd -m -p '$6$MtwVl8yJ/PeTevqf$C.7ANQ7EzLLP9b6dXte9keVRYeUDDZ2i579yEIA4nPRpx8DMvxkRx3mZSYbZHCdzh8rmvq/PAvV9zTwp/i6o1/' -g sudo -G sudo tester
COPY . .
RUN rm *.tar
RUN pip3 install -r ./requirements.txt
RUN apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
EXPOSE 22 6000 7000 8000 9000 30303
CMD ["/usr/sbin/sshd", "-D"]
CMD ["/lib/systemd/systemd"]
