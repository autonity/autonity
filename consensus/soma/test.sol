pragma solidity ^0.4.23;

contract Test {
  function test() public pure returns(string) {
    return "Hello Test!!!";
  }

  int private count = 0;
  function incrementCounter() public {
    count += 1;
  }
  function decrementCounter() public {
    count -= 1;
  }
  function getCount() public view returns (int) {
    return count;
  }
}
