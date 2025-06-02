// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity >=0.8.2 < 0.9.0;

import {ReentrancyGuard} from "./ReentrancyGuard.sol";
import "./interfaces/IAutonity.sol";
import {IConfigEvents} from "./interfaces/IConfigEvents.sol";
import "./interfaces/IOracle.sol";
import {EnumerableSet} from "./utils/Set.sol";

/**
 * @title Autonity Protocol - Oracle Contract
 * @notice This contract implements the Oracle for the Autonity Protocol, allowing voters to submit price reports
 * and aggregate them while detecting outliers.
 */
contract Oracle is IOracle, IConfigEvents, ReentrancyGuard {
    using EnumerableSet for EnumerableSet.AddressSet;

    // Struct to hold metadata information concerning a voter
    struct VoterInfo {
        uint256 round; // The last round the voter participated in
        uint256 commit; // The commit hash of the voter's last report
        uint256 performance; // The performance score of the voter
        uint256 nonRevealCount; // Number of commits that were not revealed
        bool isVoter; // Indicates if the address is a registered voter
        bool reportAvailable; // Indicates if the last report is available for the voter
    }

    // Struct to hold price data
    struct Price {
        uint256 price; // The reported price
        uint timestamp; // The timestamp of the price report
        bool success; // Indicates if the price report was successful
    }

    // Configuration settings for the Oracle
    struct Config {
        IAutonity autonity; // Address of the Autonity contract
        address operator; // Address of the operator
        uint votePeriod; // Duration of the voting period
        int256 outlierDetectionThreshold; // Threshold for outlier detection
        int256 outlierSlashingThreshold; // Threshold for slashing outliers
        uint256 baseSlashingRate; // Base rate for slashing
        uint256 nonRevealThreshold; // Threshold for missed reveals
        uint256 revealResetInterval; // Number of rounds after the missed reveal counter is reset
        uint256 slashingRateCap; // maximum slashing rate for oracle penalties
    }

    // ==== Public state variables ====
    Config internal config;
    int256 internal symbolUpdatedRound = type(int256).min; // todo(youssef): confused why not uint256
    uint256 internal lastRoundBlock;
    mapping(address => VoterInfo) internal voterInfo;
    mapping(string => mapping(address => Report)) internal reports;

    // ==== Private state variables ====
    // @dev Note that the oracle DECIMALS cannot be changed without having an effect on the
    // Stabilization computations
    string[] private symbols;
    string[] private newSymbols;
    uint private newVotePeriod;

    address[] private voters;
    address[] private newVoters;
    bool private newVotersSet;
    bool private newVotersAccessUpdated;
    uint256 private round;
    mapping(string => Price)[] internal prices;

    // voting periods are not necessarily in sync with epoch rewards, so we need to keep track of
    // rewards for each voter over the epoch, and distribute for all voting periods that concluded this epoch
    EnumerableSet.AddressSet private rewardReceivers;

    // in the voterInfo mapping we only store the performance for the current voting round
    // in order to persist information between epochs and voter sets, we keep a separate accumulating
    // mapping for all rounds ended this epoch
    mapping(address => address) internal voterTreasuries;
    mapping(address => address) internal voterValidators;
    mapping(address => uint256) private rewardPeriodPerformance;
    uint256 private rewardPeriodAggregatedScore;

    /**
     * @dev Constructor to initialize the Oracle contract.
     * @param _voters List of initial voters' addresses.
     * @param _symbols List of symbols to be tracked.
     * @param _config Configuration settings for the oracle.
     */
    constructor(
        address[] memory _voters,
        address[] memory _nodeAddresses,
        address[] memory _treasuries,
        string[] memory _symbols,
        Config memory _config
    ) {
        require(
            _config.nonRevealThreshold < _config.revealResetInterval && _config.revealResetInterval > 0,
            "invalid config"
        );
        config = _config;
        symbols = _symbols;
        newSymbols = _symbols;

        for (uint i = 0; i < _voters.length; i++) {
            voterTreasuries[_voters[i]] = _treasuries[i];
            voterValidators[_voters[i]] = _nodeAddresses[i];
            voterInfo[_voters[i]].isVoter = true;
        }

        _votersSort(_voters, int(0), int(_voters.length - 1));
        voters = _voters;
        newVoters = _voters;
        newVotePeriod = config.votePeriod;
        round = 1;
        // create the space for first index in prices array
        prices.push();
        _checkVotePeriod(_config.votePeriod);
    }

    /**
    * @dev Receive Auton function
    * https://solidity.readthedocs.io/en/v0.7.2/contracts.html#receive-ether-function
    */
    receive() external payable {}

    /**
    * @dev Fallback function
    * https://solidity.readthedocs.io/en/v0.7.2/contracts.html#fallback-function
    */
    fallback() external payable {}

    /**
     * @notice Allows a validator to vote for the current period using a commit-reveal scheme.
     *
     * This function implements a commit-reveal voting mechanism for validators in the consensus committee.
     * Validators can only vote once per round, and their votes will be discarded if they leave the
     * consensus committee or not counted if they try to vote again in the same round.
     *
     * If a validator attempts to vote multiple times in a round, the transaction will revert.
     * Important: Non-reverting calls are refunded by the Autonity Protocol.
     *
     * @param _commit A hash representing the commitment of the new reports.
     * @param _reports An array of `Report` objects revealing the reports for the previous cycle.
     * @param _salt The salt value used to generate the last round's commitment.
     * @param _extra May be used for carrying Autonity Server version. Not used internally.
     *
     * @dev
     * - The function checks that the validator has not already voted in the current round.
     * - If the voter is new (first round), their vote will not be recorded until their next submission.
     * - If the number of reports does not match the expected symbols, the vote will be discarded.
     * - If the last round's commitment does not match the current commitment, the vote is invalidated.
     *
     */
    function vote(
        uint256 _commit,
        Report[] calldata _reports,
        uint256 _salt,
        uint8 _extra
    ) onlyVoters external virtual nonReentrant {
        // revert if already voted for this round
        // voters should not be allowed to vote multiple times in a round
        // because we are refunding the tx fee and this opens up the possibility
        // to spam the node
        VoterInfo storage _voterInfo = voterInfo[msg.sender];
        require(_voterInfo.round != round, "already voted");

        uint256 _pastCommit = _voterInfo.commit;
        // Store the new commit before checking against reveal to ensure an updated commit is
        // available for the next round in case of failures.
        _voterInfo.commit = _commit;
        uint256 _lastVotedRound = _voterInfo.round;
        // considered to be voted whether vote is valid or not
        _voterInfo.round = round;
        // new voter/first round
        if (_lastVotedRound == 0) {
            emit NewVoter(msg.sender, _extra);
            return;
        }

        if (_lastVotedRound != round - 1) {
            emit InvalidVote("LastVotedRoundMismatch", msg.sender, round - 1, _lastVotedRound, _extra);
            return;
        }

        _commit = uint256(keccak256(abi.encode(_reports, _salt, msg.sender)));
        if (_pastCommit != _commit) {
            _increaseNonRevealCount(msg.sender, _voterInfo);
            emit InvalidVote("CommitMismatch", msg.sender, _pastCommit, _commit, _extra);
            // we return the tx fee in all cases, because in both cases voter is slashed during aggregation
            // phase, because the reports contain invalid prices
            return;
        }

        // if data is not supplied and voter is not a new voter
        // report must contain the correct price
        if (_reports.length != symbols.length) {
            emit InvalidVote("ReportLengthMismatch", msg.sender, _reports.length, symbols.length, _extra);
            return;
        }

        // Voter voted on every symbols
        for (uint256 i = 0; i < _reports.length; i++) {
            require(_reports[i].confidence <= 100, "invalid confidence score");
            require(
                (_reports[i].price > 0 && _reports[i].confidence > 0),
                "confidence/price error"
            );
            reports[symbols[i]][msg.sender] = _reports[i];
        }
        _voterInfo.reportAvailable = true;
        emit SuccessfulVote(msg.sender, _extra);
    }

    /**
     * @notice Finalizes the current round and aggregates the votes. Called by the Autonity contract.
     * @return true if there is a new round and new symbol prices are available, false if not.
     * @dev This function has technically infinite gas budget and must not throw in any condition.
     */
    function finalize() onlyAutonity external virtual nonReentrant returns (bool) {
        if (block.number < lastRoundBlock + config.votePeriod) {
            return false;
        }

        // first penalize for no reveal
        bool[] memory _penalizedVoters = _penalizeForNoReveal();
        if (round % config.revealResetInterval == 0) {
            _resetNonRevealCounter();
        }

        prices.push();
        for (uint i = 0; i < symbols.length; i += 1) {
            _aggregateReports(i, _penalizedVoters);
        }

        _finalizeRewards();

        lastRoundBlock = block.number;
        round += 1;
        // apply the new vote period at the end of the round to get the oracle network synced with it.
        if (config.votePeriod != newVotePeriod) {
            config.votePeriod = newVotePeriod;
        }
        emit NewRound(round, block.timestamp, config.votePeriod);
        _resetPenalizedReports(_penalizedVoters);
        return true;
    }

    function _finalizeRewards() internal virtual {
        for (uint256 i = 0; i < voters.length; i++) {
            address _voter = voters[i];
            if (voterInfo[_voter].performance > 0) {
                rewardReceivers.add(_voter);
                rewardPeriodPerformance[_voter] += voterInfo[_voter].performance;
                rewardPeriodAggregatedScore += voterInfo[_voter].performance;
                voterInfo[_voter].performance = 0;
            }
        }
    }

    /**
     * @notice Distributes oracle rewards to voters based on their performance. Called by Autonity at finalize().
     * @param _ntn, the amount of ntn to redistribute
     */
    function distributeRewards(uint256 _ntn) onlyAutonity external virtual nonReentrant payable {
        uint256 _atn = address(this).balance;
        _performRewardDistribution(_atn, _ntn);
    }

    function _performRewardDistribution(uint256 _totalATN, uint256 _totalNTN) internal virtual {
        if (rewardPeriodAggregatedScore == 0) {
            return;
        }
        address[] memory _receivers = rewardReceivers.values();
        for (uint256 i = 0; i < _receivers.length; i++) {
            address _voter = _receivers[i];
            uint256 _atn = (_totalATN * rewardPeriodPerformance[_voter]) / rewardPeriodAggregatedScore;
            uint256 _ntn = (_totalNTN * rewardPeriodPerformance[_voter]) / rewardPeriodAggregatedScore;

            // Transfer ATN rewards
            // funds for failed transfers will be redistributed for the next round
            // emit an event to notify the voter
            if(_atn > 0) {
                (bool _sent, bytes memory _returnData) = voterTreasuries[_voter].call{value: _atn, gas: 2300}("");
                if (_sent == false) {
                    emit IAutonity.CallFailed(voterTreasuries[_voter], "", _returnData);
                }
            }

            // Transfer NTN rewards
            if(_ntn > 0) {
                config.autonity.autobond(voterValidators[_voter], _ntn, 0);
            }

            rewardPeriodPerformance[_voter] = 0;
            rewardReceivers.remove(_voter);
        }
        emit TotalOracleRewards(_totalNTN, _totalATN);
        rewardPeriodAggregatedScore = 0;
    }

    /**
     * @notice updates voters and symbols, taking into account boundary edge cases.
     *         Called by Autonity at the start of a new Oracle round
     */
    function updateVotersAndSymbol() onlyAutonity external virtual {
        // this votingInfo is updated with the newVoter set just so that the new voters
        // are able to send their first vote, but they will not be used for aggregation
        // in this round
        if (newVotersSet == true) {
            for (uint i = 0; i < newVoters.length; i++) {
                voterInfo[newVoters[i]].isVoter = true;
            }
            newVotersAccessUpdated = true;
            newVotersSet = false;
        }
        else if (newVotersAccessUpdated == true) {
            // votingInfo update happens a round later then setting of new voters,
            // because we still want to aggregate vote for lastVoterSet
            _updateVotingInfo();
            newVotersAccessUpdated = false;
        }

        // symbol update should happen in the symbolUpdatedRound+2 since we expect
        // oracles to send commit for newSymbols in symbolUpdatedRound+1 and reports
        // for the new symbols in symbolUpdatedRound+2
        if (int256(round) == symbolUpdatedRound + 2) {
            symbols = newSymbols;
            for (uint i = 0; i < voters.length; i++) {
                voterInfo[voters[i]].reportAvailable = false;
            }
        }
    }

    /**
     * @notice Aggregates reports for a specific symbol.
     * @param _sindex The index of the symbol to aggregate.
     * @dev This function detects outliers and calculates the final price for the symbol.
     */
    function _aggregateReports(uint _sindex, bool[] memory _penalizedVoters) internal virtual {
        string memory _symbol = symbols[_sindex];
        Report[] memory _totalReports = new Report[](voters.length);
        uint256 _count = 0;
        for (uint i = 0; i < voters.length; i++) {
            address _voter = voters[i];
            // if there is no available report from this validator we must account for it.
            if (!voterInfo[_voter].reportAvailable) {
                continue;
            }
            _totalReports[_count++] = reports[_symbol][_voter];
        }

        // at this stage if count > 0 we must have valid strictly positive reports available.
        if (_count > 0) {
            int256 _priceMedian = int256(uint256(_getMedian(_totalReports, _count)));
            // exclude and detect outliers
            OutlierDetection memory _outlierDetection = _findOutliers(_priceMedian, _symbol);

            if (_outlierDetection.totalReports > 0) {
                // punish outliers if found
                for (uint256 i = 0; i < _outlierDetection.totalOutliers; i++) {
                    uint256 _slashingAmount = _penalize(
                        _outlierDetection.outliers[i],
                        _priceMedian,
                        reports[_symbol][voters[_outlierDetection.outliers[i]]],
                        _penalizedVoters
                    );
                    emit Penalized(
                        voters[_outlierDetection.outliers[i]],
                        _slashingAmount,
                        _symbol,
                        _priceMedian,
                        reports[_symbol][voters[_outlierDetection.outliers[i]]].price
                    );
                }
                uint256 _price = _calculateWeightedPrice(
                    _outlierDetection.filteredReports,
                    _outlierDetection.totalReports
                );
                prices[round][_symbol] = Price(
                    _price,
                    block.timestamp,
                    true
                );
                emit PriceUpdated(_price, round, _symbol, true, block.timestamp);
            } else {
                uint256 _price = prices[round - 1][_symbol].price;
                // all voters are detected as outliers, so no valid report found
                // use past value for price if unsuccesful
                prices[round][_symbol] = Price(
                    _price,
                    block.timestamp,
                    false
                );
                emit PriceUpdated(_price, round, _symbol, false, block.timestamp);
            }
        } else {
            uint256 _price = prices[round - 1][_symbol].price;
            // use past value for price if unsuccesful
            prices[round][_symbol] = Price(
                _price,
                block.timestamp,
                false
            );
            emit PriceUpdated(_price, round, _symbol, false, block.timestamp);
        }
    }

    /**
     * @notice Return latest available price data.
     * @param _symbol, the symbol from which the current price should be returned.
     */
    function latestRoundData(string memory _symbol) external view virtual nonReentrantView returns (RoundData memory data) {
        //return last aggregated round
        Price memory _p = prices[round - 1][_symbol];
        RoundData memory _d = RoundData(round - 1, _p.price, _p.timestamp, _p.success);
        return _d;
    }

    /**
     * @notice Return price data for a specific round.
     * @param _round, the round for which the price should be returned.
     * @param _symbol, the symbol for which the current price should be returned.
     * @dev IOracle interface method
     */
    function getRoundData(uint256 _round, string memory _symbol) external view virtual nonReentrantView returns (RoundData memory data) {
        Price memory _p = prices[_round][_symbol];
        RoundData memory _d = RoundData(_round, _p.price, _p.timestamp, _p.success);
        return _d;
    }

    /**
     * @return config, the current oracle config
     */
    function getConfig() external virtual view nonReentrantView returns (Config memory){
        return config;
    }

    /**
     * @param _oracleAddress, the oracle address of a validator
     * @return his node address
     */
    function getVoterValidators(address _oracleAddress) external virtual view nonReentrantView returns (address) {
        return voterValidators[_oracleAddress];
    }

    /**
     * @param _oracleAddress, the oracle address of a validator
     * @return his treasury address
     */
    function getVoterTreasuries(address _oracleAddress) external virtual view nonReentrantView returns (address) {
        return voterTreasuries[_oracleAddress];
    }

    /**
    * @return the round at which the symbols got updated
    */
    function getSymbolUpdatedRound() external virtual view returns (int256){
        return symbolUpdatedRound;
    }

    /**
    * @return the block at which the last completed round ended
    */
    function getLastRoundBlock() external virtual view nonReentrantView returns (uint256){
        return lastRoundBlock;
    }

    /**
    * @param _voter, the voter address
    * @return the related voter information
    */
    function getVoterInfo(address _voter) external virtual view nonReentrantView returns (VoterInfo memory){
        return voterInfo[_voter];
    }

    /**
    * @param _symbol, the target symbol
    * @param _voter, the target voter address
    * @return the latest report of that voter for that symbol
    */
    function getReports(string memory _symbol, address _voter) external virtual view nonReentrantView returns (Report memory){
        return reports[_symbol][_voter];
    }

    /**
     * @notice Update the symbols to be requested.
     * Only effective at the next round.
     * Restricted to the operator account.
     * @param _symbols list of string symbols to be used. E.g. "ATN-USD"
     * @dev emit {NewSymbols} event.
     * @dev IOracle interface method
     */
    function setSymbols(string[] memory _symbols) external virtual onlyOperator {
        require(_symbols.length != 0, "symbols can't be empty");
        require((symbolUpdatedRound + 1 != int256(round)) && (symbolUpdatedRound != int256(round)), "can't be updated in this round");
        newSymbols = _symbols;
        symbolUpdatedRound = int256(round);
        // these symbols will be effective for oracles from next round
        emit NewSymbols(_symbols, round + 1);
    }

    /**
     * @notice Retrieve the lists of symbols to be voted on.
     */
    function getSymbols() external virtual view returns (string[] memory) {
        // if current round is the next round of the symbol update round
        // we should return the updated symbols, because oracle clients are supposed
        // to use updated symbols to fetch data
        if (symbolUpdatedRound + 1 == int256(round)) {
            return newSymbols;
        }
        return symbols;
    }

    /**
     * @notice Retrieve the list of new participants in the Oracle process.
     * @dev IOracle interface method implementation.
     */
    function getNewVoters() external virtual view nonReentrantView returns (address[] memory) {
        return newVoters;
    }

    /**
     * @notice Retrieve the list of participants in the Oracle process.
     * @dev IOracle interface method implementation.
     */
    function getVoters() external virtual view nonReentrantView returns (address[] memory) {
        return voters;
    }

    /**
     * @notice Returns the tolerance for missed reveal count before the voter gets punished.
     */
    function getNonRevealThreshold() external virtual view returns (uint256) {
        return config.nonRevealThreshold;
    }

    /**
    * @notice Retrieve the current round ID.
    * @dev IOracle interface method implementation.
    */
    function getRound() external virtual view nonReentrantView returns (uint256) {
        return round;
    }

    /**
    * @notice Decimal places to be used with price reports
    * @dev IOracle interface method implementation.
    */
    function getDecimals() external pure returns (uint8) {
        return DECIMALS;
    }

    /**
    * @notice Retrieve the performance for a voter in this reward (epoch) period.
    */
    function getRewardPeriodPerformance(address _voter) external virtual view nonReentrantView returns (uint256) {
        return rewardPeriodPerformance[_voter];
    }

    /**
    * @notice Retrieve the vote period.
    * @dev IOracle interface method implementation.
    */
    function getVotePeriod() external virtual view nonReentrantView returns (uint) {
        return config.votePeriod;
    }

    /**
    * @notice Retrieve the new vote period that is going to be applied at the end of the vote round.
    * @dev IOracle interface method implementation.
    */
    function getNewVotePeriod() external virtual view nonReentrantView returns (uint) {
        return newVotePeriod;
    }

    /**
     * @notice Called to update the list of the oracle voters.
     * @dev Only accessible from the Autonity Contract.
     * @dev IOracle interface method implementation.
     */
    function setVoters(
        address[] memory _newVoters,
        address[] memory _treasury,
        address[] memory _validator
    ) onlyAutonity external virtual nonReentrant {
        require(_newVoters.length != 0, "Voters can't be empty");
        for (uint256 i = 0; i < _newVoters.length; i++) {
            voterTreasuries[_newVoters[i]] = _treasury[i];
            voterValidators[_newVoters[i]] = _validator[i];
        }
        _votersSort(_newVoters, int(0), int(_newVoters.length - 1));
        newVoters = _newVoters;
        newVotersSet = true;
    }

    /**
    * @notice Setter for the operator.
    * @dev IOracle interface method implementation.
    */
    function setOperator(address _operator) external virtual onlyAutonity {
        config.operator = _operator;
    }

    /**
    * @notice Setter for the vote period, new vote period will be applied at the end of the round.
    * @dev IOracle interface method implementation..
    */
    function setVotePeriod(uint _votePeriod) external virtual onlyOperator {
        _checkVotePeriod(_votePeriod);
        newVotePeriod = _votePeriod;
        emit ConfigUpdateUint("votePeriod", config.votePeriod, _votePeriod, lastRoundBlock + config.votePeriod);
    }

    /**
     * @notice Setter for commit-reveal penalty mechanism configuration.
     */
    function setCommitRevealConfig(uint256 _threshold, uint256 _resetInterval) external virtual onlyOperator {
        require(
            _threshold < _resetInterval && _resetInterval > 0,
            "invalid config"
        );
        emit ConfigUpdateUint("revealResetInterval", config.revealResetInterval, _resetInterval, block.number);
        config.revealResetInterval = _resetInterval;
        emit ConfigUpdateUint("nonRevealThreshold", config.nonRevealThreshold, _threshold, block.number);
        config.nonRevealThreshold = _threshold;
    }

    /**
    * @notice Setter for the internal slashing and outlier detection configuration.
    */
    function setSlashingConfig(
        int256 _outlierSlashingThreshold,
        int256 _outlierDetectionThreshold,
        uint256 _baseSlashingRate,
        uint256 _slashingRateCap
    ) external virtual onlyOperator {
        emit ConfigUpdateInt("outlierSlashingThreshold", config.outlierSlashingThreshold, _outlierSlashingThreshold, block.number);
        config.outlierSlashingThreshold = _outlierSlashingThreshold;
        emit ConfigUpdateInt("outlierDetectionThreshold", config.outlierDetectionThreshold, _outlierDetectionThreshold, block.number);
        config.outlierDetectionThreshold = _outlierDetectionThreshold;
        emit ConfigUpdateUint("baseSlashingRate", config.baseSlashingRate, _baseSlashingRate, block.number);
        config.baseSlashingRate = _baseSlashingRate;
        emit ConfigUpdateUint("slashingRateCap", config.slashingRateCap, _slashingRateCap, block.number);
        config.slashingRateCap = _slashingRateCap;
    }

    function _checkVotePeriod(uint _votePeriod) internal virtual view {
        // we need this check to update new voters at the end of voting round
        uint256 _epochPeriod = config.autonity.getCurrentEpochPeriod();
        require(_votePeriod * 2 <= _epochPeriod, "vote period is too big");
        _epochPeriod = config.autonity.getEpochPeriod();
        require(_votePeriod * 2 <= _epochPeriod, "vote period is too big");
    }

    function _updateVotingInfo() internal virtual {
        uint _i = 0;
        uint _j = 0;

        while (_i < voters.length && _j < newVoters.length) {
            if (voters[_i] == newVoters[_j]) {
                _i++;
                _j++;
                continue;
            } else if (voters[_i] < newVoters[_j]) {
                // delete from votingInfo since this voter is not present in the new Voters
                delete voterInfo[voters[_i]];
                _i++;
            } else {
                _j++;
            }
        }

        while (_i < voters.length) {
            // delete from voted since it's not present in the new Voters
            delete voterInfo[voters[_i]];
            _i++;
        }
        voters = newVoters;
    }

    /**
    * @dev QuickSort algorithm sorting addresses in lexicographic order.
    */
    function _votersSort(address[] memory _voters, int _low, int _high) internal virtual pure {
        if (_low >= _high) return;
        int _i = _low;
        int _j = _high;
        address _pivot = _voters[uint(_low + (_high - _low) / 2)];
        while (_i <= _j) {
            while (_voters[uint(_i)] < _pivot) _i++;
            while (_voters[uint(_j)] > _pivot) _j--;
            if (_i <= _j) {
                (_voters[uint(_i)], _voters[uint(_j)]) =
                (_voters[uint(_j)], _voters[uint(_i)]);
                _i++;
                _j--;
            }
        }
        // Recursion call in the left partition of the array
        if (_low < _j) {
            _votersSort(_voters, _low, _j);
        }
        // Recursion call in the right partition
        if (_i < _high) {
            _votersSort(_voters, _i, _high);
        }
    }

    /**
    * @dev QuickSort algorithm sorting addresses in lexicographic order.
    */
    function _getMedian(Report[] memory _priceArray, uint _length) internal virtual pure returns (uint120) {
        if (_length == 0) {
            return 0;
        }
        _sortPrice(_priceArray, 0, int(_length - 1));
        uint _midIndex = _length / 2;
        return (_length % 2 == 0) ?
            (_priceArray[_midIndex - 1].price + _priceArray[_midIndex].price) / 2 : _priceArray[_midIndex].price;
    }

    function _sortPrice(Report[] memory _priceArray, int _low, int _high) internal virtual pure {
        int _i = _low;
        int _j = _high;
        if (_i == _j) return;
        uint120 pivot = _priceArray[uint(_low + (_high - _low) / 2)].price;
        while (_i <= _j) {
            while (_priceArray[uint(_i)].price < pivot) _i++;
            while (pivot < _priceArray[uint(_j)].price) _j--;
            if (_i <= _j) {
                (_priceArray[uint(_i)], _priceArray[uint(_j)]) = (_priceArray[uint(_j)], _priceArray[uint(_i)]);
                _j--;
                _i++;
            }
        }
        // recurse left partition
        if (_low < _j) {
            _sortPrice(_priceArray, _low, _j);
        }
        // recurse right partition
        if (_i < _high) {
            _sortPrice(_priceArray, _i, _high);
        }
        return;
    }

    /// @dev Return struct for outlier detection
    struct OutlierDetection {
        uint256[] outliers;
        uint256 totalOutliers;
        Report[] filteredReports;
        uint256 totalReports;
    }

    /**
     * @dev Identifies and returns outlier voter addresses based on their price reports relative to a given median.
     *
     * This internal function processes a list of voters and checks their reported prices against the provided median.
     * It categorizes voters into outliers and non-outliers based on a configurable detection threshold.
     * The function also accumulates non-outlier reports and updates the performance of each voter.
     *
     * Note: This function iterates through the state, which may result in high gas costs.
     * Consider optimizing for gas savings in future implementations.
     *
     * @param _median The median price to compare against for outlier detection.
     * @param _symbol The symbol representing the price report to analyze.
     *
     * @return OutlierDetection Struct containing:
     *         - outliers: Array of outlier indices.
     *         - totalOutliers: Total number of outliers detected.
     *         - filteredReports: Array of reports that are not outliers.
     *         - totalReports: Total number of reports that are not outliers.
     */
    function _findOutliers(int256 _median, string memory _symbol) internal virtual returns (OutlierDetection memory){
        OutlierDetection memory result;

        result.filteredReports = new Report[](voters.length);
        result.outliers = new uint256[](voters.length);

        for (uint256 i = 0; i < voters.length; i++) {
            address _voter = voters[i];
            if (!voterInfo[_voter].reportAvailable) {
                continue;
            }
            // median here is assumed to be non-0.
            // we don't want the following to underflow
            int256 _ratio = (_median - int256(uint256(reports[_symbol][_voter].price))) * 100 / _median;
            if (_ratio <= config.outlierDetectionThreshold && - 1 * _ratio <= config.outlierDetectionThreshold) {
                result.filteredReports[result.totalReports++] = reports[_symbol][_voter];
                // take advantage of this iteration to include performance calculation
                voterInfo[_voter].performance += reports[_symbol][_voter].confidence;
            } else {
                result.outliers[result.totalOutliers++] = i;
            }
        }
        return result;
    }

    function _calculateWeightedPrice(Report[] memory _report, uint256 _reportCount) internal virtual pure returns (uint256) {
        uint256 _totalConfidence = 0;
        uint256 _price = 0;
        for (uint256 i = 0; i < _reportCount; i++) {
            _price += _report[i].price * _report[i].confidence;
            _totalConfidence += _report[i].confidence;
        }
        return _price / _totalConfidence;
    }

    function _penalize(uint256 _outlierIndex, int256 _median, Report memory _report, bool[] memory _penalizeVoters) internal virtual returns (uint256) {
        address _outlier = voters[_outlierIndex];
        // Stop considering this reporter for any future calculation.
        // This is symbol independant.
        voterInfo[_outlier].reportAvailable = false;
        if (_penalizeVoters[_outlierIndex]) {
            // already penalized for no reveal
            return 0;
        }
        int256 _diffRatio = (int256(uint256(_report.price)) - _median) * 100 / _median;
        //price is 120 bits max so _diffratio squared is at most 240 bits
        _diffRatio = _diffRatio * _diffRatio;
        if (_diffRatio <= config.outlierSlashingThreshold) {
            return 0;
        }

        // `_diffRatio` is a percentage squared, so dividing it by 10_000
        uint256 _slashingRate = (uint256(_diffRatio - config.outlierSlashingThreshold) *
                               uint256(_report.confidence) *
                               config.baseSlashingRate) / 10_000;

        // Capped the oracle slashing rate
        if (_slashingRate > config.slashingRateCap) {
            _slashingRate = config.slashingRateCap;
        }

        return config.autonity.slash(voterValidators[_outlier], _slashingRate);
    }

    function _penalizeForNoReveal() internal virtual returns (bool[] memory) {
        bool[] memory penalizedVoters = new bool[](voters.length);
        // penalize for commit without reveal
        for (uint i = 0; i < voters.length; i++) {
            address _voter = voters[i];
            VoterInfo storage _voterInfo = voterInfo[_voter];
            /*
                Following two scenarios can happen:
                    1. Voter voted in this round but did not reveal his commit from
                        the last round (it's not the first round for the voter).
                    2. Voter provided commit in the last round but did not vote in this round.
                In both cases, the voter submitted commit in the last round but did not reveal.
                Point 2 is handled here and point 1 is handled in `vote` function.
            */
            if (_voterInfo.round == round - 1 && _voterInfo.commit > 0) {
                _increaseNonRevealCount(_voter, _voterInfo);
            }

            if (_voterInfo.nonRevealCount > config.nonRevealThreshold) {
                penalizedVoters[i] = true;
                emit NoRevealPenalty(_voter, round, _voterInfo.nonRevealCount);
                _voterInfo.nonRevealCount = 0;
                // penalize with highest
                config.autonity.slash(_voter, config.slashingRateCap);
            }
        }
        return penalizedVoters;
    }

    function _resetNonRevealCounter() internal virtual {
        for (uint i = 0; i < voters.length; i++) {
            voterInfo[voters[i]].nonRevealCount = 0;
        }
    }

    function _resetPenalizedReports(bool[] memory penalizedVoters) internal virtual {
        for (uint i = 0; i < voters.length; i++) {
            if (penalizedVoters[i]) {
                voterInfo[voters[i]].reportAvailable = false;
            }
        }
    }

    function _increaseNonRevealCount(address _address, VoterInfo storage _voter) internal virtual {
        _voter.nonRevealCount++;
        emit CommitRevealMissed(_address, round, _voter.nonRevealCount);
    }

    /*
     ============================================================
         Modifiers
     ============================================================
     */
    modifier onlyVoters {
        require(voterInfo[msg.sender].isVoter, "restricted to only voters");
        _;
    }

    modifier onlyAutonity {
        require(address(config.autonity) == msg.sender, "restricted to the autonity contract");
        _;
    }

    modifier onlyOperator {
        require(config.operator == msg.sender, "restricted to operator");
        _;
    }
}
