// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity >=0.8.2 < 0.9.0;
/**
 * @dev Interface of the Oracle Contract
 */
interface IOracle {

    struct RoundData {
        uint256 round;
        int256 price;
        uint timestamp;
        uint status;
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
     *
     * @dev Emit a {Vote} event in case of succesful vote.
     *
     * @param _commit hash of the ABI packed-encoded prevotes to be
     * submitted the next voting round.
     * @param _reports list of prices to be voted on. Ordering must
     * respect the list of symbols returned by {getSymbols}.
     *
     */
    function vote(uint256 _commit, int256[] memory _reports, uint256 _salt ) external;
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
     * @notice Called to update the list of the oracle voters.
     * @dev Only accessible from the Autonity Contract.
     */
    function setVoters(address[] memory _newVoters) external;

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
     * @notice Retrieve the current round ID.
    */
    function getRound() external view returns (uint256);
    /**
    * @notice Precision to be used with price reports
    */
    function getPrecision() external view returns (uint256);


    /**
     * @dev Emitted when a vote has been succesfully accounted after a {vote} call.
     */
    event Voted(address indexed _voter, int[] _votes);
    /**
     * @dev Emitted when a vote has been succesfully accounted after a {vote} call.
     * round - the round at which new symbols are effective
     */
    event NewSymbols(string[] _symbols, uint256 _round);
    /**
     * @dev Emitted when a new voting round is started.
     * round - the new round ID
     * height - the height of the current block being executed in the EVM context.
     * timestamp - the TS in time's seconds since Jan 1 1970 (Unix time) that the block been mined by protocol
     * votePeriod - the round period in blocks for the price voting and aggregation.
     */
    event NewRound(uint256 _round, uint256 _height, uint256 _timestamp, uint _votePeriod);
}