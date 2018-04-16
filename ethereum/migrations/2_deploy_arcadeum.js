var Arcadeum = artifacts.require("./Arcadeum.sol");

module.exports = function(deployer) {
  deployer.deploy(Arcadeum, {
    gasPrice: 1,
    gasLimit: 4000000000000000,
    gas: 4000000000000000
  });
};
