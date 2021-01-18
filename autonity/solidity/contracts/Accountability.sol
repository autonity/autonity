// SPDX-License-Identifier: MIT

pragma solidity ^0.7.1;

// how to write and use precompiled contracts https://blog.qtum.org/precompiled-contracts-and-confidential-assets-55f2b47b231d
library Accountability {

    // a proof is used to prove a suspicious or an innocent behavior.
    struct Proof {
        // <rule, height, round, msgType, sender> form the identity for the on-chain management of proof.
        uint rule;          // rule id of accountability patterns.
        uint height;
        uint round;
        uint msgType;
        address sender;
        bytes message;      // raw bytes of the message to be proved.
        bytes[][] evidence; // raw bytes of the messages as evidence of a behavior.
    }

    // use for bytes decoding in EVM when calls precompiled contract
    struct StaticProof {
        uint rule;
        uint height;
        uint round;
        uint msgType;
        address sender;
        uint numOfMessages;   // save number of messages in the byte array.
        uint[] lengthOfEach; // save number of bytes for each message.
        bytes messages;       // bytes array for all the msgs, the 1st slot is the one to be suspicious.
    }

    // call precompiled contract to check if challenge is valid, for the node who is on the challenge
    // need to issue a proof of innocent via transaction to get the challenge to be removed on a reasonable
    // time window which is a system configuration.
    function takeChallenge(Proof memory challenge) internal view returns (uint[2] memory p) {
        // todo: assemble static structure and calculate the byte array to be copied into EVM context.
        StaticProof memory cProof;
        uint len = 0;

        assembly {
        //staticcall(gasLimit, to, inputOffset, inputSize, outputOffset, outputSize)
            if iszero(staticcall(gas(), 0xfd, cProof, len, p, 0x40)) {
                revert(0, 0)
            }
        }

        return p;
    }

    // call precompiled contract to check if proof of innocent is valid or not, the caller will remove
    // the challenge if the proof is valid.
    function innocentCheck(Proof memory innocent) internal view returns (uint[2] memory p) {
        // todo: assemble static structure and calculate the byte array to be copied into EVM context.
        StaticProof memory iProof;
        uint len = 0;

        assembly {
        //staticcall(gasLimit, to, inputOffset, inputSize, outputOffset, outputSize)
            if iszero(staticcall(gas(), 0xfe, iProof, len, p, 0x40)) {
                revert(0, 0)
            }
        }

        return p;
    }
}
