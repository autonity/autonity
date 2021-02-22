// SPDX-License-Identifier: MIT

pragma solidity ^0.7.1;

// how to write and use precompiled contracts https://blog.qtum.org/precompiled-contracts-and-confidential-assets-55f2b47b231d
library Accountability {
    struct Proof {
        address sender;
        bytes32 msghash;
        // the rlp encoded Proof. Please check afd_types.go type RaWProof struct.
        bytes rawproof;
    }

    // checkChallenge, it validate proof of challenge is valid by according afd rules, the precompiled contract returns
    // the msg sender address and the msg hash when the proof is valid.
    function checkChallenge(bytes memory proof) internal view returns (address, bytes32) {
        uint length = proof.length;
        uint256[2] memory retVal;
        assembly {
        //staticcall(gasLimit, to, inputOffset, inputSize, outputOffset, outputSize)
            if iszero(staticcall(gas(), 0xfd, proof, length, retVal, 64)) {
                revert(0, 0)
            }
        }

        return (address(retVal[0]), bytes32(retVal[1]));
    }

    // checkAccusation, it validate the accusation is valid by according rules and it returns the msg sender address
    // and the msg hash when the accusation is valid.
    function checkAccusation(bytes memory proof) internal view returns (address, bytes32) {
        uint length = proof.length;
        uint256[2] memory retVal;

        assembly {
        //staticcall(gasLimit, to, inputOffset, inputSize, outputOffset, outputSize)
            if iszero(staticcall(gas(), 0xfc, proof, length, retVal, 64)) {
                revert(0, 0)
            }
        }

        return (address(retVal[0]), bytes32(retVal[1]));
    }

    // checkInnocent, it validate the proof of innocent is valid by according to afd rules, the precompiled contract returns
    // the msg sender address and the msg hash when the proof is valid.
    function checkInnocent(bytes memory proof) internal view returns (address, bytes32) {
        uint length = proof.length;
        uint256[2] memory retVal;

        assembly {
        //staticcall(gasLimit, to, inputOffset, inputSize, outputOffset, outputSize)
            if iszero(staticcall(gas(), 0xfe, proof, length, retVal, 64)) {
                revert(0, 0)
            }
        }

        return (address(retVal[0]), bytes32(retVal[1]));
    }
}
