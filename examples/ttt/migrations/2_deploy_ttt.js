var Arcadeum = artifacts.require('arcadeum-contracts/Arcadeum');
var TTT = artifacts.require('./TTT');

module.exports = function(deployer) {
  deployer.deploy(TTT, Arcadeum.address);
};
