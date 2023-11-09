// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;

contract MockActivityKeyProofVerifier {

    constructor(){}

    fallback(bytes calldata input) external payable returns (bytes memory) {
        uint256 ret = 1;
        return abi.encodePacked(ret);
    }
}
