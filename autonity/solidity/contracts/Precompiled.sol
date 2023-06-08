// SPDX-License-Identifier: MIT

pragma solidity ^0.8.3;

// how to write and use precompiled contracts https://blog.qtum.org/precompiled-contracts-and-confidential-assets-55f2b47b231d
library Precompiled {
    address constant INNOCENCE_CONTRACT = address(0xfd);
    address constant ACCUSATION_CONTRACT = address(0xfc);
    address constant MISBEHAVIOUR_CONTRACT = address(0xfe);
    function enodeCheck(string memory _enode) internal view returns (address, uint) {
        uint[2] memory p;
        address addr;
        assembly {
            //staticcall(gasLimit, to, inputOffset, inputSize, outputOffset, outputSize)
            if iszero(staticcall(gas(), 0xff, add(_enode,32), mload(_enode), p, 0x40)) {
                revert(0, 0)
            }
            addr :=  div(mload(p), 0x1000000000000000000000000) // abi encoded, shift >> 96
        }
        return (addr, p[1]);
    }

    // call precompiled contracts with its corresponding contract address and the rlp encoded accountability event.
    // return a tuple that contains the corresponding address of the validator, the consensus msg's hash and the
    // verification result of the corresponding accountability event, the rule id of the event and the corresponding
    // height of the accountability event against to.
    // returns(msgSender, msgHash, result, ruleID)
    function checkAccountabilityEvent(address _to, bytes memory _proof) internal view returns
        (address, bytes32, uint256, uint256) {
        // type bytes in solidity consumes the first 32 bytes to save the length of the byte array, thus the memory copy
        // in the static call should take the extra 32 bytes to have all the rlp encoded bytes copied, otherwise the
        // decoding of rlp would fail.
        uint length = _proof.length + 32;
        uint256[4] memory retVal;
        assembly {
            //staticcall(gasLimit, to, inputOffset, inputSize, outputOffset, outputSize)
            if iszero(staticcall(gas(), _to, _proof, length, retVal, 128)) {
                revert(0, 0)
            }
        }

        return (address(uint160(retVal[0])), bytes32(retVal[1]), retVal[2], retVal[3]);
    }
}
