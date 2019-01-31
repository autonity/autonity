pragma solidity ^0.4.23;
pragma experimental ABIEncoderV2;


/* Soma.sol
*
*  Soma is a smart contract implementation of Clique consensus, see EIP 225,
*  however validators no longer vote for proposers to be added or removed through the Soma contract.
*/

contract Soma {

    address[] public validators;


    constructor (address[] _validators) public {

        for (uint256 i = 0; i < _validators.length; i++) {
            validators.push(_validators[i]);
        }

    }


    function AddValidator(address _validator) public onlyValidators(msg.sender) {
        //Just need to make sure we're duplicating the entry
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

    function getValidators() public view returns (address[]) {
        return validators;
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