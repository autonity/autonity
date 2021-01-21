// SPDX-License-Identifier: MIT

pragma solidity ^0.7.1;

// how to write and use precompiled contracts https://blog.qtum.org/precompiled-contracts-and-confidential-assets-55f2b47b231d
library Accountability {
    // Proof is used for precompiled contract, since solidity's limitation and
    // the performance concern, data packing was done by golang in client side, then AC
    // just bypass the data to precompiled contracts.
    // https://docs.soliditylang.org/en/v0.8.0/internals/layout_in_storage.html
    // https://docs.soliditylang.org/en/v0.8.0/internals/layout_in_memory.html
    struct Proof {
        uint256 height;       // 32 bytes
        address sender;       // 20 bytes
        uint64 round;         // 8 bytes
        uint8 rule;           // 1 bytes
        uint8 msgType;        // 1 bytes

        // A dynamic byte array which contains formatted meta data for precompiled contract to unpack proof, the pack
        // and unpack is done is golang client side since concern on performance and Solidity's limitations on data packing.
        bytes packedProof;
    }

    // call precompiled contract to check if challenge is valid
    function checkChallenge(bytes memory proof) internal view returns (uint[2] memory p) {
        uint length = proof.length;

        assembly {
        //staticcall(gasLimit, to, inputOffset, inputSize, outputOffset, outputSize)
            if iszero(staticcall(gas(), 0xfd, proof, length, p, 0x40)) {
                revert(0, 0)
            }
        }

        return p;
    }

    // call precompiled contract to check if proof of innocent is valid or not, the caller will remove
    // the challenge if the proof is valid.
    function checkInnocent(bytes memory proof) internal view returns (uint[2] memory p) {
        uint length = proof.length;

        assembly {
        //staticcall(gasLimit, to, inputOffset, inputSize, outputOffset, outputSize)
            if iszero(staticcall(gas(), 0xfe, proof, length, p, 0x40)) {
                revert(0, 0)
            }
        }

        return p;
    }
}
