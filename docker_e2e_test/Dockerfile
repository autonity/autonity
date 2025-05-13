FROM ubuntu:24.04

ENV PYTHON_VERSION=3.8 \
    DEBIAN_FRONTEND=noninteractive \
    VENV=/opt/venv

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        software-properties-common \
        build-essential \
        wget \
        zlib1g-dev \
        libncurses5-dev \
        libgdbm-dev \
        libnss3-dev \
        libssl-dev \
        libreadline-dev \
        libffi-dev \
        libsqlite3-dev \
        openssh-client && \
    add-apt-repository ppa:deadsnakes/ppa -y && \
    apt-get install -y --no-install-recommends \
        python${PYTHON_VERSION} \
        python${PYTHON_VERSION}-venv \
        python${PYTHON_VERSION}-dev && \
    apt-get clean && \
    rm -f /usr/bin/python3 && \
    ln -s /usr/bin/python${PYTHON_VERSION} /usr/bin/python3 && \
    ln -s /usr/bin/python${PYTHON_VERSION} /usr/bin/python

RUN python3 -m venv ${VENV}
ENV PATH="${VENV}/bin:$PATH"

RUN pip install --no-cache-dir --upgrade pip setuptools

RUN python -c "import pkg_resources; print('setuptools version:', pkg_resources.get_distribution('setuptools').version)"

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .
EXPOSE 22 6000 7000