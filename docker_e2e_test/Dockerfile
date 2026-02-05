# Docker file for the test engine of the autonity docker e2e test framework.
FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive \
    VENV=/opt/venv

# Use distro Python (Ubuntu 24.04 ships Python 3.12). The engine dependencies
# are kept in a dedicated requirements file that is compatible with this Python.
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        python3 \
        python3-venv \
        python3-pip \
        ca-certificates \
        openssh-client && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN python3 -m venv ${VENV}
ENV PATH="${VENV}/bin:$PATH"

RUN python3 -m pip install --no-cache-dir --upgrade pip setuptools wheel

COPY requirements_engine.txt .
RUN python3 -m pip install --no-cache-dir -r requirements_engine.txt

COPY . .
EXPOSE 22 6000 7000
