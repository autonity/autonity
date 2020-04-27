pragma solidity ^0.6.4;

// how to write and use precompiled contracts https://blog.qtum.org/precompiled-contracts-and-confidential-assets-55f2b47b231d
library Precompiled {
    function enodeCheck(string memory _enode) internal view returns (uint[2] memory p) {
        assembly {
            //staticcall(gasLimit, to, inputOffset, inputSize, outputOffset, outputSize)
            if iszero(staticcall(gas(), 0xff, _enode, 0xc0, p, 0x40)) {
                revert(0, 0)
            }
        }

        return p;
    }
}
