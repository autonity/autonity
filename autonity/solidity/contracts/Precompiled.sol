// SPDX-License-Identifier: MIT

pragma solidity ^0.7.1;

/*
// how to write and use precompiled contracts https://blog.qtum.org/precompiled-contracts-and-confidential-assets-55f2b47b231d
library Precompiled {
    function enodeCheck(string memory _enode) internal view returns (uint[2] memory p) {
        // cast string to bytes does works in some case, for example string is utf8 encoded
        // that is why deployment of autonity contract is reverted by evm.
        uint calldata_len = bytes(_enode).length;
        assembly {
            //staticcall(gasLimit, to, inputOffset, inputSize, outputOffset, outputSize)
            if iszero(staticcall(gas(), 0xff, _enode, calldata_len, p, 0x40)) {
                revert(0, 0)
            }
        }
        return p;
    }
}*/

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
}
