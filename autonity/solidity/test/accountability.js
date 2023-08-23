/* TODO(tariq) tests for accountability.sol
 *
 * high priority:
 * 1. test that verifies the stake slashing logic (_slash function)
 *    1. total slash = (_slashingRate * _availableFunds)/config.slashingRatePrecision;
 *    2. first self unbonding stake gets slashed (PAS)
 *    3. then self bonded stake gets slashed (PAS)
 *    4. remaining slash gets split between delegated stake and delegated unbonding stake
 *    5. slashed validator gets jailed for the amount of time we expect
 *    6. slashed funds gets moved to the autonity treasury (this particular condition might be better to address it in protocol.js, since it is an interaction with Autonity.sol... to decide)
 * 2. verify that when having multiple slashing events in the slashing queue, the offenceCount goes up (and consequently the slashing rate). _performSlashingTasks function
 * 3. verify that a validator with an history of past offences gets slashed more than a clean one (and exactly the amount more we expect).
 * 4. Accusation flow tests (test canAccuse/canSlash when appriopriate)
 *    0. issue multiple accusations on different blocks --> check that only the one who expired get converted to misbehavior
      1. Validator is accused and submit proof of innocence before the window is expired --> no slashing
      2. Validator is accused and does not submit proof of innocence --> accusation is promoted to misbehavior and validator gets slashed
      3. Validator is accused and submits proof of innocence **after** the window is expired --> accusation is promoted to misbehavior and validator gets slashed
      3. Validator is accused while already under accusation --> 2nd accusation reverts
        - canAccuse should return a deadline for when we can submit the 2nd accusation
      4. Validators is under accusation and someone sends proof of misbehavior against him (for the same epoch of the accusation). 
         The accused validator does not publish proof of innocence for the accusation. Outcome:
           - if misbehaviour severity > accusation severity --> only misbehaviour slashing takes effect
           - if misbehaviour severity < accusation severity --> both offences are slashed
  * 5. cannot submit misbehaviour for validator already slashed for the offence epoch with a higher severity than the submitted misb
  *       require(slashingHistory[_offender][_epoch] < _severity, "already slashed at the proof's epoch");
  *     - also canSlash should return false
  * 6. same thing for accusation
  *     - canAccuse should return false.
  * 7. edge scenario. validator is sentenced for 2 misbehaviour with 1st misb severity < 2nd misb severity in the same epoch. He should be slashed for both
  * 8. validator already slashed for an epoch, but accusation with higher severity is issued against him --> accusation is valid and should lead to slashing if not addressed
  * 9. other edge cases?

low priority:
 * 1 test that config gets set properly at contract deploy (low priority)
 * 2. modifiers (low priority)
 *    - only registered validators can submit accountability events (handleEvent)
 *    - only autonity can call finalize(), setEpochPeriod() and distributeRewards()
 * 3. verify rule --> severity mapping
 * 4. verify severity --> slashign rate mapping
 * 5. test chunked event processing (handleEvent function)
 *      - test also case where multiple validators are sending interleaved chunks
 * 6. test _handle* functions edge cases (e.g. invalid proof, block in future, etc.) --> tx should revert
 * 7. whitebox testing (better to leave for when the implementation will be less prone to changes)
 *    - verify that the accusation queue, the slashing queue update and the other internal structures are updated as we expect
 */
