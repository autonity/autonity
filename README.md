## Autonity

[![Join the chat at https://gitter.im/clearmatics/autonity](https://badges.gitter.im/clearmatics/autonity.svg)](https://gitter.im/clearmatics/autonity?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Build Status](https://travis-ci.org/clearmatics/autonity.svg?branch=master)](https://travis-ci.org/clearmatics/autonity)
[![Coverage Status](https://coveralls.io/repos/github/clearmatics/autonity/badge.svg?branch=tendermint)](https://coveralls.io/github/clearmatics/autonity?branch=tendermint)

Autonity is a generalization of the Ethereum protocol based on a fork of go-ethereum.

[Autonity Documentation](https://docs.autonity.io)

## Prerequisites

* Go (version 1.9 or later) - https://golang.org/dl
* A C compiler.
* Docker

## Working with the source

Before working with the source you will need to run

```
make embed-autonity-contract
```

This generates go source from the autonity contract.

## Building Autonity

```
make autonity
```

## Build Autonity docker image

```
make build-docker-image
```

## License

The go-ethereum library (i.e. all code outside of the `cmd` directory) is licensed under the
[GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html),
also included in our repository in the `COPYING.LESSER` file.

The go-ethereum binaries (i.e. all code inside of the `cmd` directory) is licensed under the
[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also
included in our repository in the `COPYING` file.
