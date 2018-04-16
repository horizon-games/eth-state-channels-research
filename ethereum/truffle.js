var HDWalletProvider = require("truffle-hdwallet-provider");

module.exports = {
  // See <http://truffleframework.com/docs/advanced/configuration>
  // to customize your Truffle configuration!
  networks: {
    ganache: {
      host: "127.0.0.1",
      port: 8545,
      network_id: "*",
      gasPrice: 1,
      gasLimit: 4000000000000000,
      gas: 4000000000000000
    },
    rinkeby: {
      provider: function () {
        return new HDWalletProvider("salon oval sausage day year song edge december tortoise elephant search review model civil wonder", "https://rinkeby.infura.io/P8djn1ELvrq7uw7LrE22")
      },
      network_id: "*",
      gasPrice: 200000000000
    }
  }
};
