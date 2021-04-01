// SPDX-License-Identifier: MIT

pragma solidity ^0.7.1;

// how to write and use precompiled contracts https://blog.qtum.org/precompiled-contracts-and-confidential-assets-55f2b47b231d
library Precompiled {
    function enodeCheck(string memory _enode) internal view returns (uint[2] memory p) {
        assembly {
            //staticcall(gasLimit, to, inputOffset, inputSize, outputOffset, outputSize)
            if iszero(staticcall(gas(), 0xff, _enode, 0xc0, p, 64)) {
                revert(0, 0)
            }
        }
        return p;
    }
    // checkChallenge, it validate proof of challenge is valid by according faultdetector rules, the precompiled contract returns
    // the msg sender address and the msg hash when the proof is valid.
    function checkMisbehaviour(bytes memory proof) internal view returns (address, bytes32, uint256) {
        // bytes in solidity consumes the first 32 bytes as the length of the byte array, so the length for memory copy should plus
        // extra 32 bytes to copy the full byte array, otherwise it would be truncated with the tail-32-bytes of the array.
        uint length = proof.length + 32;
        uint256[3] memory retVal;
        assembly {
            //staticcall(gasLimit, to, inputOffset, inputSize, outputOffset, outputSize)
            if iszero(staticcall(gas(), 0xfe, proof, length, retVal, 96)) {
                revert(0, 0)
            }
        }

        return (address(retVal[0]), bytes32(retVal[1]), retVal[2]);
    }

    // checkAccusation, it validate the accusation is valid by according rules and it returns the msg sender address
    // and the msg hash when the accusation is valid.
    function checkAccusation(bytes memory proof) internal view returns (address, bytes32, uint256) {
        // bytes in solidity consumes the first 32 bytes as the length of the byte array, so the length for memory copy should plus
        // extra 32 bytes to copy the full byte array, otherwise it would be truncated with the tail-32-bytes of the array.
        uint length = proof.length + 32;
        uint256[3] memory retVal;

        assembly {
            //staticcall(gasLimit, to, inputOffset, inputSize, outputOffset, outputSize)
            if iszero(staticcall(gas(), 0xfc, proof, length, retVal, 96)) {
                revert(0, 0)
            }
        }

        return (address(retVal[0]), bytes32(retVal[1]), retVal[2]);
    }

    // checkInnocence, it validate the proof of innocent is valid by according to faultdetector rules, the precompiled contract returns
    // the msg sender address and the msg hash when the proof is valid.
    function checkInnocence(bytes memory proof) internal view returns (address, bytes32, uint256) {
        // bytes in solidity consumes the first 32 bytes as the length of the byte array, so the length for memory copy should plus
        // extra 32 bytes to copy the full byte array, otherwise it would be truncated with the tail-32-bytes of the array.
        uint length = proof.length + 32;
        uint256[3] memory retVal;

        assembly {
            //staticcall(gasLimit, to, inputOffset, inputSize, outputOffset, outputSize)
            if iszero(staticcall(gas(), 0xfd, proof, length, retVal, 96)) {
                revert(0, 0)
            }
        }

        return (address(retVal[0]), bytes32(retVal[1]), retVal[2]);
    }
}
