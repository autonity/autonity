pragma solidity ^0.5.0;

library Precompiled {
    function enodeCheck(string enode) internal pure returns (uint[2] memory p) {
        assembly {
            if iszero(staticcall(gas, 0xff, input, 0x80, p, 0x40)) {
                revert(0, 0)
            }
        }

        return p;
    }
}