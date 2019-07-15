pragma solidity ^0.5.1;
pragma experimental ABIEncoderV2;

contract Autonity {

    // list of validators of network
    address[] public validators;
    // enodesWhitelist - which nodes can connect to network
    string[] public enodesWhitelist;
    address public owner;



    // constructor get called at block #1 with msg.owner equal to Soma's deployer
    // configured in the genesis file.
    constructor (address[] memory _validators, string[] memory _enodesWhitelist) public {
        for (uint256 i = 0; i < _validators.length; i++) {
            validators.push(_validators[i]);
        }

        for (uint256 i = 0; i < _enodesWhitelist.length; i++) {
            enodesWhitelist.push(_enodesWhitelist[i]);
        }
        owner = msg.sender;
    }


    function AddValidator(address _validator) public onlyValidators(msg.sender) {
        //Need to make sure we're duplicating the entry
        validators.push(_validator);
    }


    function RemoveValidator(address _validator) public onlyValidators(msg.sender) {

        require(validators.length > 1);

        for (uint256 i = 0; i < validators.length; i++) {
            if (validators[i] == _validator){
                validators[i] = validators[validators.length - 1];
                validators.length--;
                break;
            }
        }

    }

    function AddEnode(string memory  _enode) public {
        //Need to make sure we're not duplicating the entry
        enodesWhitelist.push(_enode);
    }


    function RemoveEnode(string memory  _enode) public {

        require(enodesWhitelist.length > 1);

        for (uint256 i = 0; i < enodesWhitelist.length; i++) {
            if (compareStringsbyBytes(enodesWhitelist[i], _enode)) {
                enodesWhitelist[i] = enodesWhitelist[enodesWhitelist.length - 1];
                enodesWhitelist.length--;
                break;
            }
        }

    }

    function compareStringsbyBytes(string memory s1, string memory s2) public pure returns(bool){
        return keccak256(abi.encodePacked(s1)) == keccak256(abi.encodePacked(s2));
    }

    /*
    ========================================================================================================================

        Getters - extra values we may wish to return

    ========================================================================================================================
    */

    /*
    * getValidators
    *
    * Returns the macro validator list
    */

    function GetValidators() public view returns (address[] memory) {
        return validators;
    }

    /*
    * getWhitelist
    *
    * Returns the macro participants list
    */

    function getWhitelist() public view returns (string[] memory) {
        return enodesWhitelist;
    }


    /*
    ========================================================================================================================

        Modifiers

    ========================================================================================================================
    */

    /*
    * onlyValidators
    *
    * Modifier that checks if the voter is an active validator
    */

    modifier onlyValidators(address _voter) {
        bool present = false;
        for (uint256 i = 0; i < validators.length; i++) {
            if(validators[i] == _voter){
                present = true;
                break;
            }
        }
        require(present, "Voter is not a validator");
        _;
    }

}