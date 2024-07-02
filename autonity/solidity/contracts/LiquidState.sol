// SPDX-License-Identifier: LGPL-3.0-only

pragma solidity ^0.8.3;
import "./LiquidLogic.sol";
import "./Autonity.sol";

contract LiquidState is Upgradeable {

    // storage layout - this must be compatible with LiquidLogic
    mapping(address => uint256) private balances;
    mapping(address => uint256) private lockedBalances;

    mapping(address => mapping (address => uint256)) private allowances;
    uint256 private supply;

    mapping(address => uint256) private atnRealisedFees;
    mapping(address => uint256) private atnUnrealisedFeeFactors;
    uint256 private atnLastUnrealisedFeeFactor;
    mapping(address => uint256) private ntnRealisedFees;
    mapping(address => uint256) private ntnUnrealisedFeeFactors;
    uint256 private ntnLastUnrealisedFeeFactor;

    string public name;
    string public symbol;

    address public validator;
    address payable public treasury;
    uint256 public commissionRate;

    uint256 public treasuryUnclaimedATN;

    Autonity private autonityContract; //not hardcoded for testing purposes


    constructor(
        address _validator,
        address payable _treasury,
        uint256 _commissionRate,
        string memory _index,
        address _liquidLogicAddress
    ) {
        // commissionRate <= 1.0
        require(_commissionRate <= LiquidLogic(payable(_liquidLogicAddress)).COMMISSION_RATE_PRECISION());

        validator = _validator;
        treasury = _treasury;
        commissionRate = _commissionRate;
        name = string.concat("LNTN-", _index);
        symbol = string.concat("LNTN-", _index);
        autonityContract = Autonity(payable(msg.sender));
    }

    /**
     * @dev Fallback function that delegates calls to the address returned by `liquidLogicContract()`. Will run if no other
     * function in the contract matches the call data.
     */
    fallback() payable external {
        _delegate(
            liquidLogicContract()
        );
    }

    /**
     * @dev Fallback function that delegates calls to the address returned by `liquidLogicContract()`. Will run if call data
     * is empty.
     */
    receive() payable external {
        _delegate(
            liquidLogicContract()
        );
    }

    /**
     ============================================================

        Internals

     ============================================================
     */

    /**
     * @notice Fetch liquid logic contract address from autonity
     */
    function liquidLogicContract() public view returns (address) {
        address _address = autonityContract.liquidLogicContract();
        require(_address != address(0), "liquid logic contract not set");
        return _address;
    }

    /**
     * @dev Do a static call to `_contractAddress`. Use for pure functions.
     */
    function _static(address _contractAddress, bytes memory _data) internal view returns (bytes memory) {
        (bool _success, bytes memory _returnData) = _contractAddress.staticcall(_data);
        require(_success, "static call failed");
        return _returnData;
    }

    /**
     * @dev Delegates the current call to `_contractAddress`.
     * 
     * This function does not return to its internall call site, it will return directly to the external caller.
     */
    function _delegate(address _contractAddress) internal {
        // solhint-disable-next-line no-inline-assembly
        assembly {
            // Copy msg.data. We take full control of memory in this inline assembly
            // block because it will not return to Solidity code. We overwrite the
            // Solidity scratch pad at memory position 0.
            calldatacopy(0, 0, calldatasize())

            // Call the implementation.
            // out and outsize are 0 because we don't know the size yet.
            let result := delegatecall(gas(), _contractAddress, 0, calldatasize(), 0, 0)

            // Copy the returned data.
            returndatacopy(0, 0, returndatasize())

            if iszero(result) {
                revert(0, returndatasize())
            }
            return(0, returndatasize())
        }
    }

    /*
     ============================================================
        Getters
     ============================================================
     */

    /**
     * @notice Returns the total claimable fees (AUT) earned by the delegator to-date.
     * @dev The logic to compute claimable rewards are in LiquidLogic.sol. A static call is done to LiquidLogic contract
     * to fetch the claimable rewards. The function unclaimedRewards() in LiquidLogic is a pure function so that we can
     * make a static call and compute the rewards. Otherwise it will view state from LiquidLogic contract instead of
     * Liquidstate contract. And we cannot do a delegate call from a view function. So we need unclaimedRewards() in
     * LiquidLogic to be a pure function and make a static call.
     * @param _account the delegator account.
     */
    function unclaimedRewards(address _account) external view returns(uint256 _unclaimedATN, uint256 _unclaimedNTN) {
        bytes memory _returnData = _static(
            liquidLogicContract(),
            abi.encodeWithSignature(
                "unclaimedRewards(uint256,uint256,uint256,uint256,uint256,uint256,uint256)",
                balances[_account], atnRealisedFees[_account], ntnRealisedFees[_account],
                atnUnrealisedFeeFactors[_account], ntnUnrealisedFeeFactors[_account],
                atnLastUnrealisedFeeFactor, ntnLastUnrealisedFeeFactor
            )
        );

        require(_returnData.length >= 64, "not enough return data");

        bytes memory _unclaimedATNBytes = new bytes(32);
        bytes memory _unclaimedNTNBytes = new bytes(32);

        assembly {
            mstore(add(_unclaimedATNBytes, 32), mload(add(_returnData, 32)))
            mstore(add(_unclaimedNTNBytes, 32), mload(add(_returnData, 64)))
        }

        _unclaimedATN = uint256(bytes32(_unclaimedATNBytes));
        _unclaimedNTN = uint256(bytes32(_unclaimedNTNBytes));
    }

    /**
     * @notice Returns the total amount of stake token issued.
     */
    function totalSupply() external view returns (uint256) {
        return supply;
    }

    /**
     * @return uint8 the number of decimals the LNTN token uses.
     * @dev ERC-20 Optional.
     */
    function decimals() public pure returns (uint8) {
        return DECIMALS;
    }

    /**
     * @notice Returns the amount of liquid newtons held by the account (ERC-20).
     */
    function balanceOf(address _delegator) external view returns (uint256) {
        return balances[_delegator];
    }

    /**
     * @notice Returns the amount of locked liquid newtons held by the account.
     */
    function lockedBalanceOf(address _delegator) external view returns (uint256) {
        return lockedBalances[_delegator];
    }

    /**
     * @notice Returns the amount of unlocked liquid newtons held by the account.
     */
    function unlockedBalanceOf(address _delegator) external view returns (uint256) {
        return  balances[_delegator] - lockedBalances[_delegator];
    }

    /**
     * @dev See {IERC20-allowance}.
     */
    function allowance(address _owner, address _spender) public view returns (uint256) {
        return allowances[_owner][_spender];
    }
}
