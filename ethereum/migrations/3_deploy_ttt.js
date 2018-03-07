var Arcadeum = artifacts.require("./Arcadeum.sol");
var TTT = artifacts.require("./TTT.sol");

module.exports = function(deployer) {
  deployer.deploy(TTT, Arcadeum.address);
};
