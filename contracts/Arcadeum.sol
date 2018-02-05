pragma solidity ^0.4.19;

import './DGame.sol';

contract Arcadeum {
  mapping(address => uint) public balance;

  function deposit() external payable {
    balance[msg.sender] += msg.value;
  }

  function withdraw() external {
    uint value;

    require(balance[msg.sender] > 0);

    value = balance[msg.sender];
    delete balance[msg.sender];
    msg.sender.transfer(value);
  }
}
