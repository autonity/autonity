// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;

uint256 constant SLASHING_RATE_DECIMALS = 4;
uint256 constant SLASHING_RATE_SCALE_FACTOR = 10 ** SLASHING_RATE_DECIMALS;
// Todo: resolve a better cap for the protocol.
uint256 constant ORACLE_SLASHING_RATE_CAP = 1_000; // Capped with 10% slashing rate for oracle accountability.
