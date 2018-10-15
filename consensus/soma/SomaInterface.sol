pragma solidity ^0.4.23;

interface SomaInterface {  
    function ActiveValidator(address _validator) public view returns (bool);
}