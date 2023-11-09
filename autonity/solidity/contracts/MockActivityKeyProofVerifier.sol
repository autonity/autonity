// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;

contract MockActivityKeyProofVerifier {

    constructor(){}

    fallback(bytes calldata input) external payable returns (bytes memory ret) {
        ret = new bytes(32);
        ret[0] = bytes1(uint8(1));
        return ret;
    }
}
