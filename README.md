## Autonity

[![Join the chat at https://gitter.im/clearmatics/autonity](https://badges.gitter.im/clearmatics/autonity.svg)](https://gitter.im/clearmatics/autonity?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Build Status](https://travis-ci.org/clearmatics/autonity.svg?branch=master)](https://travis-ci.org/clearmatics/autonity)
[![Coverage Status](https://coveralls.io/repos/github/clearmatics/autonity/badge.svg?branch=tendermint)](https://coveralls.io/github/clearmatics/autonity?branch=tendermint)

Autonity is based on a fork of go-ethereum, that extends the Ethereum blockchain structure and protocol with delegated
Proof of Stake consensus.

Key technical features of the Autonity Protocol are:

- immediate and deterministic transaction finality in a public environment where participant nodes can join the network
  dynamically without permission
- delegated Proof of Stake consensus for committee selection and blockchain management, using the deterministic
  Tendermint BFT consensus algorithm
- dual coin tokenomics, with native coins auton and newton for utility and staking
- liquid staking for capital efficiency, staked newton yielding transferrable liquid newton redeemed for newton on
  unbonding stake.

Core technology is the Autonity Go Client (AGC), a fork of the Go Ethereum ‘Geth’ client. AGC is the reference
implementation of the Autonity Protocol and provides the main client software run by participant peer nodes in an
Autonity network system.

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

For information on connecting an Autonity node to a network see the documentation
website [here](https://musical-chainsaw-80f50d3e.pages.github.io/howto/run-aut/#connect-the-client-to-an-autonity-network)
.

## Open a javascript console to a node

The Autonity NodeJS Console is distributed as part of the Autonity Go Client Release in the `nodejsconsole` sub
directory. For users who only require a console, it is also available as a standalone binary from
the [Autonity Releases Archive](https://github.com/autonity/autonity/releases)
- `nodejsconsole-linux-amd64-<RELEASE_VERSION>.tar.gz`.

To run the Console and connect to a node, specify WebSockets as the transport and the IP address and port 8546 of the
Autonity client node you will connect to. Use WS for a local node and WSS for secure connection to a public node on an
Autonity network. For example, to connect to a node running on local host:

```
./nodeconsole/console ws://localhost:8546
```

The console is run with the `--experimental-repl-await` flag which means that you can use await from the console prompt.

E.G:

```
> await autonity.getMinimumBaseFee().call()
'5000000'
```

## License

The go-ethereum library (i.e. all code outside of the `cmd` directory) is licensed under the
[GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html), also included in our repository
in the `COPYING.LESSER` file.

The go-ethereum binaries (i.e. all code inside of the `cmd` directory) is licensed under the
[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also included in our repository in
the `COPYING` file.
