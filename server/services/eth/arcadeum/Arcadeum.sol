pragma solidity ^0.4.19;
pragma experimental ABIEncoderV2;

import './DGame.sol';

contract Arcadeum {
  mapping(address => uint) public balance;

  function Arcadeum() {
    address ian = 0xa5B06b0FF4FBF5D8C5e56F4a6783d28AF72a9a0d;
    balance[ian] = 20;
  }

  // Pure function to check if stopWithdrawal() will work; used to save on gas cost. Call
  // this before calling stopWithdrawal().
  function canStopWithdrawal(DGame game, uint32 matchID, uint timestamp, uint8 timestampV, bytes32 timestampR, bytes32 timestampS, uint8 subkeyV, bytes32 subkeyR, bytes32 subkeyS) public view returns (bool) {
    return false;
  }

  function isWithdrawing(address account) public view returns (bool) {
    return false;
  }

  // Could stop withdrawal in the future. Call this occasionally to refresh your cache.
  function couldStopWithdrawal(DGame game, uint32 matchID, uint timestamp, uint8 timestampV, bytes32 timestampR, bytes32 timestampS, uint8 subkeyV, bytes32 subkeyR, bytes32 subkeyS) public view returns (bool) {
    return false;
  }

  // Attempt to slash player withdrawing funds
  function stopWithdrawal(DGame game, uint32 matchID, uint timestamp, uint8 timestampV, bytes32 timestampR, bytes32 timestampS, uint8 subkeyV, bytes32 subkeyR, bytes32 subkeyS) public {
  }

  event withdrawalStarted(address indexed account);

  // Returns address of signer of subkey
  function subkeyParent(address subkey, uint8 subkeyV, bytes32 subkeyR, bytes32 subkeyS) public pure returns (address) {
    return 0xa5B06b0FF4FBF5D8C5e56F4a6783d28AF72a9a0d;
  }

  // Returns the player account address used to sign the timestamp signature (via the subkey)
  function playerAccount(DGame game, uint32 matchID, uint timestamp, uint8 timestampV, bytes32 timestampR, bytes32 timestampS, uint8 subkeyV, bytes32 subkeyR, bytes32 subkeyS) public pure returns (address) {
    address addr = 0xa5B06b0FF4FBF5D8C5e56F4a6783d28AF72a9a0d;
    return addr;
  }

  // Returns the signature used to return to each player on match start
  function matchHash(DGame game, uint32 matchID, uint timestamp, address[2] accounts, bytes[2] publicSeeds) public pure returns (bytes32) {
    return keccak256(address(game), matchID, timestamp, accounts[0], accounts[1], publicSeeds[0], publicSeeds[1]);
  }
}
