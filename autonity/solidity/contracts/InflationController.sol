// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import {SD59x18, sd, convert} from "./lib/prb-math-4.0.1/SD59x18.sol";
import {IInflationController} from "./interfaces/IInflationController.sol";


contract InflationController is IInflationController {
    struct Params {
        SD59x18 iInit;
        SD59x18 iTrans;
        SD59x18 aE;
        SD59x18 T;
    }
    Params internal params ;

    constructor(Params memory _params){
        params = _params;
    }

    /**
    * @notice Main function. Calculate NTN current supply delta.
    */
    function calculateSupplyDelta(uint256 _currentSupply, uint256 _lastEpochBlock, uint256 _currentBlock) public view returns (uint256) {
        SD59x18 _t0 = convert(int256(_lastEpochBlock));
        SD59x18 _t1 = convert(int256(_currentBlock));

        SD59x18 _lExp0 = (params.aE * _t0)/params.T;
        SD59x18 _lExp1 = (params.aE * _t1)/params.T;

        SD59x18 expTerm = params.iInit * (_t1 - _t0) +
                ((params.iInit - params.iTrans) * (_t1 - _t0) )/(params.aE.exp() - convert(1))  +
                ((params.iTrans - params.iInit) * params.T * (_lExp1.exp() - _lExp0.exp()))
                    /((params.aE.exp() - convert(1)) * params.aE);

        return uint256(convert(
                expTerm.exp() * convert(int256(_currentSupply)) -
                convert(int256(_currentSupply))
            ));
    }
}
