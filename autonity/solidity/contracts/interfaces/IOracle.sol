// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity >=0.8.2 < 0.9.0;
/**
 * @dev Interface of the Oracle Contract
 */
interface IOracle {
    /**
    * @notice RoundData.success Indicates the success of the round data retrieval. If `true`, the round data for the
    * requested symbol and round was successfully aggregated by the protocol. The caller should check the success code
    * before using the returned round data. If `false`, the protocol could not provide data for the requested symbol
    * and round, either because the symbol is invalid or because no data was collected for that round.
    */
    struct RoundData {
        uint256 round;
        uint256 price;
        uint timestamp;
        bool success;
    }

    struct Report {
        // uint120 can hold values up to approximately 1.3x10^36
        uint120 price;
        uint8 confidence;
    }

    /**
     * @notice Update the symbols to be requested.
     * Only effective at the next round.
     * Restricted to the operator account.
     * @dev emit {NewSymbols} event.
     */
    function setSymbols(string[] memory _symbols) external;

    /**
     * @notice Retrieve the lists of symbols to be voted on.
     * Need to be called by the Oracle Server as part of the init.
     */
    function getSymbols() external view returns(string[] memory _symbols);

    /**
     * @notice Vote for the prices with a commit-reveal scheme.
     * @dev Emit a {Vote} event in case of succesful vote.
     * @param _commit hash of the ABI packed-encoded prevotes to be
     * submitted the next voting round.
     * @param _reports list of prices to be voted on. Ordering must
     * respect the list of symbols returned by {getSymbols}.
     *
     */
    function vote(uint256 _commit,  Report[] calldata _reports, uint256 _salt, uint8 _extra) external;

    /**
     * @notice Get data about a specific round, using the roundId.
     */
    function getRoundData(uint256 _round, string memory _symbol) external
    view returns (RoundData memory data);

    /**
     * @notice  Get data about the last round
     */
    function latestRoundData(string memory _symbol) external view
    returns (RoundData memory data);

    /**
     * @notice Called once per VotePeriod part of the state finalisation function.
     * @dev Only accessible from the Autonity Contract.
     * @return true if there is a new round and new symbol prices are available, false if not.
     */
    function finalize() external returns (bool);

    /**
     * @notice Called when the previous round is ended. Updates the voter info for new voters.
     * @dev Only accessible from the Autonity Contract.
     */
    function updateVotersAndSymbol() external;


    /**
    * @dev Signal that rewards are available. Only accessible from the autonity contract.
    *
    */
    function distributeRewards(uint256 _ntnRewards) external payable;

    /**
     * @notice Called to update the list of the oracle voters.
     * @dev Only accessible from the Autonity Contract.
     */
    function setVoters(address[] memory _newVoters, address[] memory _treasury, address[] memory _validator) external;

    /**
     * @notice Called to update the governance operator account.
     * @dev Only accessible from the Autonity Contract.
     */
    function setOperator(address _operator) external;

    /**
    * @notice Retrieve the vote period.
    */
    function getVotePeriod() external view returns (uint);

    /**
    * @notice Retrieve the current voters in the committee.
    */
    function getVoters() external view returns(address[] memory);

    /**
    * @notice Retrieve the new voters in the committee.
    */
    function getNewVoters() external view returns(address[] memory);

    /**
     * @notice Retrieve the current round ID.
    */
    function getRound() external view returns (uint256);

    /**
    * @notice Scale to be used with price reports
    */
    function getDecimals() external view returns (uint8);

    /**
     * @notice Emitted when the oracle symbol list is updated
     * @param _symbols new symbol list
     * @param _round the round at which new symbols are effective
     */
    event NewSymbols(string[] _symbols, uint256 _round);

    /**
     * @notice Emitted when a new voting round is started.
     * @param _round the new round ID
     * @param _timestamp the TS in time's seconds since Jan 1 1970 (Unix time) that the block been mined by protocol
     * @param _votePeriod the round period in blocks for the price voting and aggregation.
     */
    event NewRound(uint256 _round,  uint256 _timestamp, uint _votePeriod);

    /**
     * @notice Emitted when oracle rewards are distributed
     * @param ntnReward total ntn rewards
     * @param atnReward total atn rewards
     */
    event TotalOracleRewards(uint256 ntnReward, uint256 atnReward);

    /**
     * @notice Emitted when an invalid report is submitted
     * @param cause cause of invalidation
     * @param reporter report submitter
     * @param expValue expected value in report
     * @param actualValue actual value in report
     */
    event InvalidVote(string cause, address indexed reporter, uint256 expValue, uint256 actualValue);

    /**
     * @notice Emitted when a valid report is accepted
     * @param reporter report submitter
     */
    event SuccessfulVote(address indexed reporter);

    /**
     * @notice Emitted when a new reporter submits a report
     * @param reporter report submitter
     */
    event NewVoter(address reporter);

    /**
     * @notice Emitted when a new price is calculated for a symbol
     * @param price new price
     * @param round round number
     * @param symbol symbol
     * @param status status of price calculation
     * @param timestamp timestamp of price
     */
    event PriceUpdated(uint256 price, uint256 round, string indexed symbol, bool status, uint256 timestamp);

    /**
     * @notice Emitted when a participant gets penalized as an outlier
     * @param _participant Oracle address of the validator
     * @param _symbol Outlier symbol.
     * @param _median Median price calculate for this symbol.
     * @param _reported Reported outlier price.
     * @param _slashingAmount Slashing amount of the validator stakes. It can be zero if the penalty does not rise above the threshold.
     */
    event Penalized(address indexed _participant, uint256 _slashingAmount, string _symbol, int256 _median, uint120 _reported);
}