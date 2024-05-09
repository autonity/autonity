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
        // Constant IR post T
        SD59x18 iPerm;
        // Note: All time related parameters MUST be denominated in seconds.
    }

    Params public params;

    constructor(Params memory _params){
        params = _params;
    }

    /**
    * @notice Main function. Calculate NTN inflation.
    */
    function calculateSupplyDelta(
        uint256 _currentSupply,
        uint256 _inflationReserve,
        uint256 _lastEpochTime,
        uint256 _currentEpochTime
    )
        external
        view
        returns (uint256)
    {
        SD59x18 _lastTime = convert(int256(_lastEpochTime));
        SD59x18 _currentTime = convert(int256(_currentEpochTime));
        if (_currentTime <= params.T) {
            return calculateTransitionRegime(_currentSupply, _lastTime, _currentTime);
        }
        if (_lastTime < params.T && _currentTime > params.T){
            uint256 _untilT = calculateTransitionRegime(_currentSupply, _lastTime, params.T);
            uint256 _afterT = calculatePermanentRegime(_inflationReserve, params.T, _currentTime);
            return _untilT + _afterT;
        }
         return calculatePermanentRegime(_inflationReserve, _lastTime, _currentTime);
    }

   /**
    * @dev Temporary. To compare against the other function for numerical precision checks.
    */
    function calculateSupplyDeltaOLD(
        uint256 _currentSupply,
        uint256 _lastEpochTime,
        uint256 _currentEpochTime
    )
        public
        view
        returns (uint256)
    {
       SD59x18 _t0 = convert(int256(_lastEpochTime));
       SD59x18 _t1 = convert(int256(_currentEpochTime));

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

    /**
    * @notice Calculate inflation before transition.
    */
    function calculateTransitionRegime(
        uint256 _currentSupply,
        SD59x18 _lastEpochTime,
        SD59x18 _currentEpochTime
    )
        internal
        view
        returns (uint256)
    {
        SD59x18 _rate;
        if (params.aE == convert(0)){
            _rate =  params.iInit + (params.iTrans - params.iInit) * _lastEpochTime / params.T;
        } else {
            SD59x18 _lExp = (params.aE * _lastEpochTime)/params.T;
            SD59x18 _rFact = (_lExp.exp() - convert(1))/ (params.aE.exp() - convert(1));
            _rate =  params.iInit + (params.iTrans - params.iInit) * _rFact;
        }
        SD59x18 _result = _rate * convert(int256(_currentSupply)) * (_currentEpochTime - _lastEpochTime);
        return uint256(convert(_result));
    }

    /**
    * @notice Calculate inflation after transition.
    */
    function calculatePermanentRegime(
        uint256 _inflationReserve,
        SD59x18 _lastEpochBlock,
        SD59x18 _currentEpochBlock
    )
        internal
        view
        returns (uint256)
    {
        return uint256(convert(
            convert(int256(_inflationReserve)) *  (_currentEpochBlock -  _lastEpochBlock) * params.iPerm )
        );
    }

}
