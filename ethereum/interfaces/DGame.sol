pragma solidity ^0.4.23;
pragma experimental ABIEncoderV2;

contract DGame {
  function matchDuration() public pure returns (uint);

  function isSecretSeedValid(address account, bytes secretSeed) public view returns (bool);

  function secretSeedRating(bytes secretSeed) public pure returns (uint32);

  function publicSeed(bytes secretSeed) public pure returns (bytes);
}
