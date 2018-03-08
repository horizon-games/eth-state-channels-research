pragma solidity ^0.4.19;
pragma experimental ABIEncoderV2;

contract DGame {
  // Constant function to check if account owns game "deck"
  function isSecretSeedValid(address account, bytes secretSeed) public view returns (bool) {
    return true;
  }

  // Get deck rating by "deck" owner
  function secretSeedRating(address account, bytes secretSeed) public view returns (uint32) {
    return 1;
  }

  // Get Merkle tree root hash of "deck"
  function publicSeed(address account, bytes secretSeed) public view returns (bytes) {
    return secretSeed;
  }
}
