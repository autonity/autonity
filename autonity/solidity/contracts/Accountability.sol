// SPDX-License-Identifier: MIT

pragma solidity ^0.7.1;

// how to write and use precompiled contracts https://blog.qtum.org/precompiled-contracts-and-confidential-assets-55f2b47b231d
library Accountability {
    struct Proof {
        // sender's address hash, used by proof reader to check if one is on challenge.
        uint256 senderHash;
        // the hash of msg which is on challenge, for distinguish auto-incriminating & equivocation msg.
        // also the precompiled contract should return the hash of the msg that was proved base on the evidence.
        // In autonity contract side, check the return hash of precompiled contract equals to the msgHash here,
        // to make sure we manage the correct proof on-chain, it also prevent byzantine node from rising proofs which was
        // accounted before.
        uint256 msgHash;

        // the rlp encoded Proof. Please check afd_types.go type RaWProof struct.
        bytes rawProof;
    }

    // call precompiled contract to check if challenge is valid, the return array is hash value of msg payload, and msg
    // sender. [PayloadHash, SenderAddressHash]
    function checkChallenge(bytes memory proof) internal view returns (uint256[2] memory p) {
        uint length = proof.length;

        assembly {
        //staticcall(gasLimit, to, inputOffset, inputSize, outputOffset, outputSize)
            if iszero(staticcall(gas(), 0xfd, proof, length, p, 0x200)) {
                revert(0, 0)
            }
        }

        return p;
    }

    // call precompiled contract to check if proof of innocent is valid or not, the caller will remove
    // the challenge if the proof is valid, the return array is hash value of msg payload, and msg sender
    // [PayloadHash, SenderAddressHash]
    function checkInnocent(bytes memory proof) internal view returns (uint256[2] memory p) {
        uint length = proof.length;

        assembly {
        //staticcall(gasLimit, to, inputOffset, inputSize, outputOffset, outputSize)
            if iszero(staticcall(gas(), 0xfe, proof, length, p, 0x200)) {
                revert(0, 0)
            }
        }

        return p;
    }
}
