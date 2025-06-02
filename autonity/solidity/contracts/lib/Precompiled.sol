// SPDX-License-Identifier: LGPL-3.0-only

pragma solidity ^0.8.3;

// how to write and use precompiled contracts https://blog.qtum.org/precompiled-contracts-and-confidential-assets-55f2b47b231d
library Precompiled {
    uint256 constant public SUCCESS = 1;

    /*
    * INNOCENCE_CONTRACT, MISBEHAVIOUR_CONTRACT and ACCUSATION_CONTRACT
    * are implemented in consensus/tendermint/accountability/contracts.go
    * All the other ones are implemented in core/vm/contracts.go
    */
    address constant public ACTIVITY_CONTRACT = address(0xf8);
    address constant public UPGRADER_CONTRACT = address(0xf9);
    address constant public COMPUTE_COMMITTEE_CONTRACT = address(0xfa);
    address constant public POP_VERIFIER_CONTRACT = address(0xfb);
    address constant public ACCUSATION_CONTRACT = address(0xfc);
    address constant public INNOCENCE_CONTRACT = address(0xfd);
    address constant public MISBEHAVIOUR_CONTRACT = address(0xfe);
    address constant public ENODE_VERIFIER_CONTRACT = address(0xff);


    function computeAbsentees(bool _mustBeEmpty, uint256 _omissionDelta, uint256 _committeeSlot) internal returns (
        bool,               // isProposerOmissionFaulty
        uint256,            // proposerEffort
        address[] memory    // absentees
    ){
        address to = ACTIVITY_CONTRACT;

        bytes memory _input = abi.encodePacked(_mustBeEmpty, _omissionDelta, _committeeSlot);
        uint256 _outputLength;

        assembly {
            // the call doesn't modify omission state, but it reads the committee from it
            // however we could actually use a static-call, it doesn't change anything from the precompile point of view
            //delegatecall(gasLimit, to, inputOffset, inputSize, outputOffset, outputSize)
            if iszero(delegatecall(gas(), to, add(_input,32), mload(_input), 0, 0)) {
                revert(0, 0)
            }
            _outputLength := returndatasize()
        }

        bytes memory _output = new bytes(_outputLength);

        assembly{
            returndatacopy(add(_output,32), 0, _outputLength)
        }

        return abi.decode(_output,(bool,uint256,address[]));
    }


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

    /*
    * @dev calls accountability precompiled contract, passing the rlp encoded proof to it
    * returns (result, offenderAddress, ruleId, msgHeight, msgHash)
    */
    function verifyAccountabilityEvent(address _to, bytes memory _proof) internal view returns
        (bool _success, address _offender, uint256 _ruleId, uint256 _block, uint256 _msgHash) {
        uint256[5] memory _returnData;
        assembly {
            //staticcall(gasLimit, to, inputOffset, inputSize, outputOffset, outputSize)
            if iszero(staticcall(gas(), _to, add(_proof,32), mload(_proof), _returnData, 160)) {
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

    /*
    * @dev calls accusation precompiled contract, passing:
    *  - the rlp encoded proof
    *  - the current values of range and delta
    * returns (result, offenderAddress, ruleId, msgHeight, msgHash)
    */
    function verifyAccountabilityAccusation(bytes memory _proof, uint256 _range, uint256 _delta, uint256 _gracePeriod) internal view returns
    (bool _success, address _offender, uint256 _ruleId, uint256 _block, uint256 _msgHash) {

        address _to = ACCUSATION_CONTRACT;
        bytes memory input = abi.encodePacked(_range, _delta, _gracePeriod ,_proof);
        uint256[5] memory _returnData;
        assembly {
        //staticcall(gasLimit, to, inputOffset, inputSize, outputOffset, outputSize)
            if iszero(staticcall(gas(), _to, add(input,32), mload(input), _returnData, 160)) {
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
