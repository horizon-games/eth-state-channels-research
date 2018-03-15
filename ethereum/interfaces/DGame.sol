pragma solidity ^0.4.19;
pragma experimental ABIEncoderV2;

contract DGame {
  // XXX: https://github.com/ethereum/solidity/issues/3270
  // *** THIS MUST MATCH Arcadeum.sol ***
  uint internal constant PUBLIC_SEED_LENGTH = 1;

  function matchDuration() public pure returns (uint);

  function isSecretSeedValid(address account, bytes secretSeed) public view returns (bool);

  function secretSeedRating(bytes secretSeed) public pure returns (uint32);

  // XXX: https://github.com/ethereum/solidity/issues/3270
  function publicSeed(bytes secretSeed) public pure returns (bytes32[PUBLIC_SEED_LENGTH]);
}
