/**
 * @dev Interface of the Autonity Contract.
 * Import this over Autonity.sol.
 */

interface IAutonity {

    /**
    * @notice Returns the current operator account.
    */
    function getOperator() external view returns (address);

    /**
    * @notice Returns the current Oracle account.
    */
    function getOracle() external view returns (address);
}
