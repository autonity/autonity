// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;

library UpdatableConfig {
    struct UintConfig {
        uint256 currentValue;
        uint256 currentActiveFrom;
        uint256 nextValue;
        uint256 nextActiveFrom;
    }

    // Returns the current value of the configuration, if the next value is active, it will return the next value
    // otherwise it will return the current value
    // @param config The configuration to read from
    function value(UintConfig storage config) internal view returns (uint256) {
        if (config.nextActiveFrom > 0 && config.nextActiveFrom <= block.timestamp) {
            return config.nextValue;
        }
        return config.currentValue;
    }

    // Updates the configuration with a new value and active from timestamp
    // @param config The configuration to update
    // @param newValue The new value to set
    // @param activeFrom The timestamp from which the new value is active
    // @return bool True if the new value overrides the next value, false otherwise
    function update(UintConfig storage config, uint256 newValue, uint256 activeFrom) internal returns (bool) {
        if (config.nextActiveFrom > 0 && config.nextActiveFrom <= block.timestamp) {
            config.currentValue = config.nextValue;
            config.nextValue = newValue;
            config.currentActiveFrom = config.nextActiveFrom;
            config.nextActiveFrom = activeFrom;
            return false;
        } else {
            bool overridden = config.nextActiveFrom > 0;
            config.nextValue = newValue;
            config.nextActiveFrom = activeFrom;
            return overridden;
        }
    }

    // Returns the current and next values of the configuration
    // @param config The configuration to read from
    // @dev notice that the current value of the config might not be the currently active value
    function current(UintConfig storage config) internal view returns (uint256, uint256) {
        return (config.currentValue, config.currentActiveFrom);
    }

    // Returns the next and next active from values of the configuration
    // @param config The configuration to read from
    // @dev notice that the next value of the config might be currently active
    function pending(UintConfig storage config) internal view returns (uint256, uint256) {
        return (config.nextValue, config.nextActiveFrom);
    }
}