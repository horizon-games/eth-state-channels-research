var TTT = artifacts.require("./TTT.sol");
var Arcadeum = artifacts.require("./Arcadeum.sol");

module.exports = function(deployer) {
  deployer.deploy(TTT, Arcadeum.address);
};
