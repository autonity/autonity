// SPDX-License-Identifier: LGPL-3.0-only

pragma solidity ^0.8.3;

// how to write and use precompiled contracts https://blog.qtum.org/precompiled-contracts-and-confidential-assets-55f2b47b231d
library Precompiled {
    uint256 constant public SUCCESS = 1;
    address constant public UPGRADER_CONTRACT = address(0xf9);
    address constant public COMPUTE_COMMITTEE_CONTRACT = address(0xfa);
    address constant public POP_VERIFIER_CONTRACT = address(0xfb);
    address constant public ACCUSATION_CONTRACT = address(0xfc);
    address constant public INNOCENCE_CONTRACT = address(0xfd);
    address constant public MISBEHAVIOUR_CONTRACT = address(0xfe);
    address constant public ENODE_VERIFIER_CONTRACT = address(0xff);

    function parseEnode(string memory _enode) internal view returns (address, uint) {
        uint[2] memory p;
        address addr;
        address to = ENODE_VERIFIER_CONTRACT;
        assembly {
            //staticcall(gasLimit, to, inputOffset, inputSize, outputOffset, outputSize)
            if iszero(staticcall(gas(), to, add(_enode,32), mload(_enode), p, 0x40)) {
                revert(0, 0)
            }
            addr :=  div(mload(p), 0x1000000000000000000000000) // abi encoded, shift >> 96
        }
        return (addr, p[1]);
    }

    /**
     * @dev Sends necessary slots to precompiled contract.
     * Committee selection and storing the committee and writing it in persistent storage are done in precompiled contract
     */
    function computeCommitteePrecompiled(uint256[5] memory input) internal {
        address to = COMPUTE_COMMITTEE_CONTRACT;
        uint256 _length = 32*5;
        assembly {
            //delegatecall(gasLimit, to, inputOffset, inputSize, outputOffset, outputSize)
            if iszero(delegatecall(gas(), to, input, _length, 0, 0)) {
                returndatacopy(0, 0, returndatasize())
                revert(0, returndatasize())
            }
        }
    }

    // call precompiled contracts with its corresponding contract address and the rlp encoded accountability event.
    // return a tuple that contains the corresponding address of the validator, the consensus msg's hash and the
    // verification result of the corresponding accountability event, the rule id of the event and the corresponding
    // height of the accountability event against to.
    // returns(msgSender, msgHash, result, ruleID, msghash)
    function verifyAccountabilityEvent(address _to, bytes memory _proof) internal view returns
        (bool _success, address _offender, uint256 _ruleId, uint256 _block, uint256 _msgHash) {
        // type bytes in solidity consumes the first 32 bytes to save the length of the byte array, thus the memory copy
        // in the static call should take the extra 32 bytes to have all the rlp encoded bytes copied, otherwise the
        // decoding of rlp would fail.
        uint _length = _proof.length + 32;
        uint256[5] memory _returnData;
        assembly {
            //staticcall(gasLimit, to, inputOffset, inputSize, outputOffset, outputSize)
            if iszero(staticcall(gas(), _to, _proof, _length, _returnData, 160)) {
                revert(0, 0)
            }
        }

        if (_returnData[0] == 1){
            _success = true;
        }
        _offender = address(uint160(_returnData[1]));
        _ruleId = _returnData[2];
        _block = _returnData[3];
        _msgHash = _returnData[4];
    }

    // @dev verify the proof of possession of validator key in a precompiled contract.
    // @param _consensusKey is a "0x" prefix hex string of the validator's BLS public key.
    // @param _proof is a "0x" prefix hex string of the proof generated together with the bls public key.
    // @param _treasury is a "0x" prefix hex string of the validator's treasury account.
    // @return 0 for a failure, 1 for a successful check.
    function popVerification(bytes memory _consensusKey, bytes memory _proof, address _treasury) internal view returns (uint256) {
        uint256[1] memory retVal;
        bytes memory input = abi.encodePacked(_consensusKey, _proof, _treasury);
        address to = POP_VERIFIER_CONTRACT;
        // type bytes in solidity consumes the first 32 bytes to save the length of the byte array, thus the memory copy
        // in the static call should take the extra 32 bytes to have all the rlp encoded bytes copied, otherwise the
        // decoding of rlp would fail.
        uint length = input.length + 32;
        assembly {
        //staticcall(gasLimit, to, inputOffset, inputSize, outputOffset, outputSize)
            if iszero(staticcall(gas(), to, input, length, retVal, 32)) {
                revert(0, 0)
            }
        }
        return retVal[0];
    }
}
