pragma solidity ^0.4.19;

import './DGame.sol';

contract Arcadeum {
  address public owner;
  mapping(address => Player) public player;

  struct Player {
    uint balance;
    uint numDeposits;
    mapping(bytes24 => uint) deposit;
  }

  function Arcadeum() public {
    owner = msg.sender;
  }

  function deposit() external payable {
    player[msg.sender].balance += msg.value;
  }

  function withdraw() external {
    uint balance;

    balance = player[msg.sender].balance;

    if (player[msg.sender].numDeposits == 0) {
      delete player[msg.sender];
    } else {
      player[msg.sender].balance = 0;
    }

    msg.sender.transfer(balance);
  }

  function lock(address game, uint32 matchID, address[] players, uint deposit) external restricted {
    bytes24 key;
    uint i;

    key = (bytes24(game) << 32) | bytes24(matchID);

    for (i = 0; i < players.length; i++) {
      player[players[i]].balance -= deposit;
      player[players[i]].numDeposits++;
      player[players[i]].deposit[key] += deposit;
    }
  }

  function unlock(address game, uint32 matchID, address[] players) external restricted {
    bytes24 key;
    uint i;
    uint deposit;

    key = (bytes24(game) << 32) | bytes24(matchID);

    for (i = 0; i < players.length; i++) {
      deposit = player[players[i]].deposit[key];
      delete player[players[i]].deposit[key];
      player[players[i]].numDeposits--;
      player[players[i]].balance += deposit;
    }
  }

  modifier restricted() {
    require(msg.sender == owner);

    _;
  }
}
