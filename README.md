## Autonity

[![Join the chat at https://gitter.im/clearmatics/autonity](https://badges.gitter.im/clearmatics/autonity.svg)](https://gitter.im/clearmatics/autonity?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Build Status](https://travis-ci.org/clearmatics/autonity.svg?branch=master)](https://travis-ci.org/clearmatics/autonity)
[![Coverage Status](https://coveralls.io/repos/github/clearmatics/autonity/badge.svg?branch=tendermint)](https://coveralls.io/github/clearmatics/autonity?branch=tendermint)

Autonity is a generalization of the Ethereum protocol based on a fork of go-ethereum.

More about the Autonity context at <https://www.autonity.io>

More detailed documentation coming soon at <https://docs.autonity.io>

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

## Building Autonity Go Client (AGC)

```
make autonity
```

## Build Autonity Go Client docker image

```
make build-docker-image
```

## Open a javascript console to a node
The address must be a websocket enabled rpc address.

```
./nodeconsole/console ws://localhost:8546
```

The console is run with the `--experimental-repl-await` flag which means that
you can use await from the console prompt.

E.G:
```
> await autonity.getMinimumGasPrice().call()
'5000'
```

## License

The go-ethereum library (i.e. all code outside of the `cmd` directory) is licensed under the
[GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html),
also included in our repository in the `COPYING.LESSER` file.

The go-ethereum binaries (i.e. all code inside of the `cmd` directory) is licensed under the
[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also
included in our repository in the `COPYING` file.
