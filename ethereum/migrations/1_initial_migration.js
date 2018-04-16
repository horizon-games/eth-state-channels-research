var Migrations = artifacts.require("./Migrations.sol");

module.exports = function(deployer) {
  deployer.deploy(Migrations, {
    gasPrice: 1,
    gasLimit: 4000000000000000,
    gas: 4000000000000000
  });
};
