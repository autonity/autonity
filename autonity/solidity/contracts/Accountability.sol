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
        address sender;     // the one who is sending a suspicious message.

        bytes message;      // raw bytes of the message to be proved.
        bytes[][] evidence; // raw bytes of the messages as evidence of a behavior.
    }

    // use for bytes decoding in EVM when calls precompiled contract
    struct ParsableProof {
        uint rule;            // 4 bytes
        uint height;          // 4 bytes
        uint round;           // 4 bytes
        uint msgType;         // 4 bytes
        address sender;       // 20 bytes
        uint numOfMessages;   // 4 bytes, save number of messages in the byte array.

        uint[] lengthOfEach;  // lengthOfEach.length * 4 bytes, save number of bytes for each message.
        bytes[] messages;     // Sum(lengthOfEach[i]) bytes, array for all the msgs, the 1st slot is the one to be suspicious.
    }

    function toParsableProof(Proof memory proof) internal returns (ParsableProof memory, uint) {
        ParsableProof memory p;
        uint totalBytes;

        p.rule = proof.rule;
        totalBytes += 4;
        p.height = proof.height;
        totalBytes += 4;
        p.round = proof.round;
        totalBytes += 4;
        p.msgType = proof.msgType;
        totalBytes += 4;
        p.sender = proof.sender;
        totalBytes += 20;
        p.numOfMessages = 1 + proof.evidence.length;
        totalBytes += 4;

        // copy the messages
        uint msgBytes = 0;
        for (uint256 i = 0; i < proof.evidence.length; i++) {
            msgBytes += proof.evidence[i].length;
            p.lengthOfEach[i] = proof.evidence[i].length;
            p.messages.push(proof.evidence[i]); // solidity does not support this operation, try packing in client side with golang.
        }

        // save the msg which is suspicious into the last slot of message set for easier decoding.
        p.lengthOfEach.push(proof.message.length);
        p.messages.push(proof.message);
        msgBytes += proof.message.length;

        totalBytes += p.lengthOfEach.length * 4;
        totalBytes +=  msgBytes;
        return (p, totalBytes);
    }

    // call precompiled contract to check if challenge is valid, for the node who is on the challenge
    // need to issue a proof of innocent via transaction to get the challenge to be removed on a reasonable
    // time window which is a system configuration.
    function takeChallenge(Proof memory challenge) internal view returns (uint[2] memory p) {
        // todo: assemble static structure and calculate the byte array to be copied into EVM context.
        ParsableProof memory cProof;
        uint len = 0;
        (cProof, len) = toParsableProof(challenge);

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
        ParsableProof memory iProof;
        uint len = 0;
        (iProof, len) = toParsableProof(innocent);

        assembly {
        //staticcall(gasLimit, to, inputOffset, inputSize, outputOffset, outputSize)
            if iszero(staticcall(gas(), 0xfe, iProof, len, p, 0x40)) {
                revert(0, 0)
            }
        }

        return p;
    }
}
