# Support setting various labels on the final image
ARG COMMIT=""
ARG VERSION=""
ARG BUILDNUM=""

# Build Autonity in a stock Go builder container
FROM golang:1.21-alpine as builder

LABEL org.opencontainers.image.source https://github.com/autonity/autonity

RUN apk add --no-cache make gcc musl-dev linux-headers libc-dev git perl-utils

ADD . /autonity
RUN cd /autonity && make autonity-docker


# Pull Autonity into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /autonity/build/bin/autonity /usr/local/bin/

EXPOSE 8545 8546 8547 30303 30303/udp
ENTRYPOINT ["autonity"]

# Add some metadata labels to help programatic image consumption
ARG COMMIT=""
ARG VERSION=""
ARG BUILDNUM=""

LABEL commit="$COMMIT" version="$VERSION" buildnum="$BUILDNUM"
