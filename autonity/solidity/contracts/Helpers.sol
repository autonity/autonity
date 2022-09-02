// SPDX-License-Identifier: MIT

pragma solidity ^0.8.3;

library Helpers {
    /**
    * @dev Converts a `uint256` to its ASCII `string` decimal representation.
     */
    function toString(uint256 value) internal pure returns (string memory) {
        // Inspired by OraclizeAPI's implementation - MIT licence
        // https://github.com/oraclize/ethereum-api/blob/b42146b063c7d6ee1358846c198246239e9360e8/oraclizeAPI_0.4.25.sol

        if (value == 0) {
            return "0";
        }
        uint256 temp = value;
        uint256 digits;
        while (temp != 0) {
            digits++;
            temp /= 10;
        }
        bytes memory buffer = new bytes(digits);
        while (value != 0) {
            digits -= 1;
            buffer[digits] = bytes1(uint8(48 + uint256(value % 10)));
            value /= 10;
        }
        return string(buffer);
    }

    /**
    * @dev Splits a `bytes` signature into v, r, s values
     */
    function splitProof(bytes memory proof) internal pure returns (bytes32 r, bytes32 s, uint8 v){
        require(proof.length == 65, "invalid proof");

        assembly {
        // signature format is packed
        // [bytes32 r] [bytes32 s] [uint8 v]
        // extract first 32 bytes
            r := mload(add(proof, 32))
        // extract second 32 bytes
            s := mload(add(proof, 64))
        // last 32 bytes,
        // we need the uint8 format and there is no uint8 mload or mload8
        // extract the byte here
            v := byte(0, mload(add(proof, 96)))
        }

        // Value for the v could possibly be [0,1,27,28] and solidity ecrecover assumes it to be 27 or 28
        // if the proof is signed with geth's crypto.sign, it sets v as [0,1]. below check updates the value
        if (v < 27) {
            v +=27;
        }
    }
}
