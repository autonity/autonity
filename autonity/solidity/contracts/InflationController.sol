// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import {SD59x18, sd, convert} from "./lib/prb-math-4.0.1/SD59x18.sol";
import {IInflationController} from "./interfaces/IInflationController.sol";


contract InflationController is IInflationController {
    struct Params {
        // Initial inflation rate
        SD59x18 iInit;
        // Transition inflation rate
        SD59x18 iTrans;
        // Convexity Parameter
        SD59x18 aE;
        // Transition Period
        SD59x18 T;

        // Note: All time related parameters MUST be denominated in seconds.
    }

    Params public params;
    uint256 public outstandingReserve;

    constructor(Params memory _params, uint256 _outstandingReserve){
        params = _params;
        outstandingReserve = _outstandingReserve;
    }

    /**
    * @notice Main function. Calculate NTN current supply delta.
    */
    function calculateSupplyDelta(uint256 _currentSupply, uint256 _inflationReserve, uint256 _lastEpochTime, uint256 _currentEpochTime) public view returns (uint256) {
        if (_currentEpochTime <= params.T) {
            return calculateTransitionRegime(_currentSupply, _lastEpochTime, _currentEpochTime);
        }
        if (_lastEpochTime < params.T && _currentEpochTime > params.T){
            uint256 _untilT = calculateTransitionRegime(_currentSupply, _lastEpochTime, params.T);
            uint256 _afterT = calculatePermanentRegime(_currentSupply, _inflationReserve, params.T, _currentEpochTime);
            return _untilT + _afterT;
        }
         return calculatePermanentRegime(_currentSupply, _inflationReserve, _lastEpochTime, _currentEpochTime);
    }

    function calculateSupplyDelta(uint256 _currentSupply,  uint256 _lastEpochTime, uint256 _currentEpochTime) public view returns (uint256) {
       SD59x18 _t0 = convert(int256(_lastEpochBlock));
       SD59x18 _t1 = convert(int256(_currentEpochBlock));

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

    function calculateTransitionRegime(uint256 _currentSupply, uint256 _lastEpochTime, uint256 _currentEpochTime) internal view returns (uint256) {
        if (params.aE == 0){

        } else {
            SD59x18 _lExp = (params.aE * _t0)/params.T;
            SD59x18 _rFact = (_lExp.exp() - convert(1))/ (params.aE.exp() - convert(1));
            SD59x18 _rate =  params.iInit + (params.iInit - params.iTrans) * _rFact;
            SD59x18 _result = _rate * convert(int256(_currentSupply)) * (_currentEpochTime - _lastEpochTime);
            return uint256(convert(result));
        }
    }

    function calculatePermanentRegime(uint256 _currentSupply, uint256 _inflationReserve, uint256 _lastEpochBlock, uint256 _currentEpochBlock) internal view returns (uint256) {

    }

}
