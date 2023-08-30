// SPDX-License-Identifier: LGPL-3.0-only

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
    * @dev extract v, r, s value from a `bytes` signature array given the starting index of signature
    * @param _multisig is the bytes array which can have one or more than one signature
    * @param _startIndex is the starting index of the signature in the byte array
     */
    function extractRSV(bytes memory _multisig, uint _startIndex) internal pure returns (bytes32 r, bytes32 s, uint8 v){
        assembly {
        // signature format is packed
        // [bytes32 r] [bytes32 s] [uint8 v]
        // extract first 32 bytes
            r := mload(add(_multisig, _startIndex))
        // extract second 32 bytes
            s := mload(add(_multisig, add(_startIndex, 32)))
        // last 32 bytes,
        // we need the uint8 format and there is no uint8 mload or mload8
        // extract the byte here
            v := byte(0, mload(add(_multisig, add(_startIndex, 64))))
        }

        // Value for the v could possibly be [0,1,27,28] and solidity ecrecover assumes it to be 27 or 28
        // if the proof is signed with geth's crypto.sign, it sets v as [0,1]. below check updates the value
        if (v < 27) {
            v +=27;
        }
    }
}
