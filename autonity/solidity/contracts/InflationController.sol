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

    uint256 internal genesisTime;

    constructor(Params memory _params){
        params = _params;
        genesisTime = block.timestamp;
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
        SD59x18 _lastTime = convert(int256(_lastEpochTime - genesisTime));
        SD59x18 _currentTime = convert(int256(_currentEpochTime - genesisTime));
        if (_currentTime <= params.T) {
            return calculateTransitionRegime(_currentSupply, _lastTime, _currentTime);
        }
        // _currentTime > params.T from here
        if (_lastTime < params.T){
            uint256 _untilT = calculateTransitionRegime(_currentSupply, _lastTime, params.T);
            uint256 _afterT = calculatePermanentRegime(_inflationReserve, params.T, _currentTime);
            return _untilT + _afterT;
        }
        return calculatePermanentRegime(_inflationReserve, _lastTime, _currentTime);
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
        SD59x18 _lastEpochTime,
        SD59x18 _currentEpochTime
    )
        internal
        view
        returns (uint256)
    {
        return uint256(convert(
            convert(int256(_inflationReserve)) *  (_currentEpochTime -  _lastEpochTime) * params.iPerm
        ));
    }

}
