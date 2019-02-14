pragma solidity ^0.4.23;
pragma experimental ABIEncoderV2;


/* Glienicke.sol
*
*  Glienicke is the network-permissioning contract.
*  The contract can be freely modified but must implement getWhitelist()
*  which should return a list with no duplicated entries of enodes
*  that are allowed to join the network.
*/

contract Glienicke {

    string[] public enodes;


    // constructor get called at block #1 with msg.owner equal to Glienicke's deployer
    // configured in the genesis file.
    constructor (string[] _genesisEnodes) public {

        for (uint256 i = 0; i < _genesisEnodes.length; i++) {
            enodes.push(_genesisEnodes[i]);
        }

    }


    function AddEnode(string _enode) public {
        //Need to make sure we're not duplicating the entry
        enodes.push(_enode);
    }


    function RemoveEnode(string _enode) public {

        require(enodes.length > 1);

        for (uint256 i = 0; i < enodes.length; i++) {
            if (compareStringsbyBytes(enodes[i], _enode)) {
                enodes[i] = enodes[enodes.length - 1];
                enodes.length--;
                break;
            }
        }

    }

    function compareStringsbyBytes(string s1, string s2) public pure returns(bool){
        return keccak256(s1) == keccak256(s2);
    }

    /*
    ========================================================================================================================

        Getters - extra values we may wish to return

    ========================================================================================================================
    */

    /*
    * getWhitelist
    *
    * Returns the macro participants list
    */

    function getWhitelist() public view returns (string[]) {
        return enodes;
    }


}