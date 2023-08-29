// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.19;

import "./Accountability.sol";

contract AccountabilityTest is Accountability {
   constructor(address payable _autonity, Config memory _config) Accountability(_autonity,_config) {}

   function slash(Event memory _event, uint256 _epochOffencesCount) public {
        Accountability._slash(_event,_epochOffencesCount);
   }
   
   function handleValidFaultProof(Event memory _event) public {
      Accountability._handleValidFaultProof(_event);
   }
   
   function handleValidAccusation(Event memory _event) public {
      Accountability._handleValidAccusation(_event);
   }

   function handleValidInnocenceProof(Event memory _event) public {
      Accountability._handleValidInnocenceProof(_event);
   }

   function performSlashingTasks() public {
      Accountability._performSlashingTasks();
   }
}
