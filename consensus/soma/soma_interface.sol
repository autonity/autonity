pragma solidity ^0.4.23;

interface SomaInterface {
    function CastVote(address _vote) public ;
    
    function ActiveValidators() public view returns (address[]);
}