Autonity Overview
=================

Autonity acts as the substrate through which a distributed group of entities
may transact in an open ended manner.

Autonity is a fork of the geth_ project modified such that it can be used by
consortiums to deploy their own private networks upon which the members are
free to deploy any contract code they wish in order to enable interaction
between themselves.

The primary distinction between geth and Autonity is that Autonity uses a proof
of stake (POS) based consensus mechanism with the actual protocol for reaching
consensus being the Tendermint_ consensus protocol. This is in contrast to
proof of work (POW) used by geth.

Tendermint is a byzantine fault tolerant (BFT) style consensus protocol, consensus is
reached via some predetermined set of participants (referred to as the
committee or inidividually as validators) coming to agreement, with each
validator's vote bearing weight equal to its stake in the system.

BFT style consensus protocols provide deterministic finality meaning that once
agreement is reached, the outcome can be considered final and will not be
subject to change. This is in contrast to POW style consensus protocols, where
one can never be entirely sure that the outcome of a transaction wont change.

Additionally Autonity provides functionality to manage the parameters of
consensus protocol. Such that the network has control over the committee
members and participants can manage their stake.

Users of an Autonity deployment pay a fee to transact through the system, the
fees are redistributed to the validators, proportional to the amount of stake
that they each have in the system.

Beyond managing the network Autonity provides no end user functionality, it is
up to the participants in the network to decide what functionality they need
and deploy it into the system.

.. _geth: https://github.com/ethereum/go-ethereum
.. _Tendermint: https://arxiv.org/pdf/1807.04938.pdf


.. toctree::
   :maxdepth: 1
   :caption: Contents

   genesis_file.rst
