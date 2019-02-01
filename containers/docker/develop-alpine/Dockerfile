FROM alpine:3.7

RUN \
  apk add --update go git make gcc musl-dev linux-headers ca-certificates && \
  git clone --depth 1 https://github.com/clearmatics/autonity && \
  (cd autonity && make autonity) && \
  cp autonity/build/bin/geth /geth && \
  apk del go git make gcc musl-dev linux-headers && \
  rm -rf /autonity && rm -rf /var/cache/apk/*

EXPOSE 8545
EXPOSE 30303

ENTRYPOINT ["/autonity"]
