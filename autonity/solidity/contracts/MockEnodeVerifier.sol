// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;

contract MockEnodeVerifier {

    constructor(){}

    // Convert an hexadecimal character to its value
    function _fromHexChar(uint8 c) internal pure returns (uint8) {
        if (bytes1(c) >= bytes1('0') && bytes1(c) <= bytes1('9')) {
            return c - uint8(bytes1('0'));
        }
        if (bytes1(c) >= bytes1('a') && bytes1(c) <= bytes1('f')) {
            return 10 + c - uint8(bytes1('a'));
        }
        if (bytes1(c) >= bytes1('A') && bytes1(c) <= bytes1('F')) {
            return 10 + c - uint8(bytes1('A'));
        }
        revert("invalid hex char");
    }
    
    // check if a char is a valid hex char
    function _validHexChar(uint8 c) internal pure returns (bool) {
      if (bytes1(c) >= bytes1('0') && bytes1(c) <= bytes1('9')) {
            return true;
        }
        if (bytes1(c) >= bytes1('a') && bytes1(c) <= bytes1('f')) {
            return true;
        }
        if (bytes1(c) >= bytes1('A') && bytes1(c) <= bytes1('F')) {
            return true;
        }
        return false;
      }

    // Convert an hexadecimal string to raw bytes
    // ref: https://ethereum.stackexchange.com/a/40247
    function _fromHex(string memory s) internal pure returns (bytes memory) {
        bytes memory ss = bytes(s);
        if (ss.length%2 != 0) { // hex string needs to be even length
            return "";
        }
        bytes memory r = new bytes(ss.length/2);
        bool c1;
        bool c2;
        for (uint i=0; i<ss.length/2; ++i) {
            c1 = _validHexChar(uint8(ss[2*i]));
            c2 = _validHexChar(uint8(ss[2*i+1]));
            if (c1 == false || c2 == false) {
            	return "";
            }
            r[i] = bytes1(_fromHexChar(uint8(ss[2*i])) * 16 + _fromHexChar(uint8(ss[2*i+1])));
        }
        return r;
    }

    fallback(bytes calldata enode) external payable returns (bytes memory ret) {
        ret = new bytes(64);
       	bytes memory err = new bytes(32); // 0 --> success, != 0 --> error
    
    	// check we have at least "enode://" + 128 bytes of pk
    	if(enode.length < 136) {
    		err[0] = bytes1(uint8(1));
    		assembly { mstore(add(ret, 64), mload(add(err,32))) }
   		    return ret;
   	    }
    
        // fetch public key encoded in hex from call data
        bytes memory publicKeyHex = new bytes(128);
        assembly { calldatacopy(add(publicKeyHex,32),add(enode.offset,8),128) }

        // convert it to raw bytes
        bytes memory publicKey = _fromHex(string(publicKeyHex));
        if (publicKey.length == 0){
   		    err[0] = bytes1(uint8(1));
    		assembly { mstore(add(ret, 64), mload(add(err,32))) }
   		    return ret;
        }

        // convert it to address
        address addr = address(uint160(uint256(keccak256(publicKey))));
        
        // pack address into bytes
        bytes memory addrBytes = new bytes(32); 
        addrBytes = abi.encodePacked(addr);

        // pack address and error into result
        assembly { mstore(add(ret, 32), mload(add(addrBytes,32))) }
        assembly { mstore(add(ret, 64), mload(add(err,32))) }
    }    
}
