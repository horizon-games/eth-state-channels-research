module.exports = {
  // See <http://truffleframework.com/docs/advanced/configuration>
  // to customize your Truffle configuration!

  // TODO: truffle has a bug preventing these non-standard paths from working
  // however Im sure in a few weeks we can try again and it'll be resolved.
  // contracts_directory: 'ethereum/contracts',
  // migrations_directory: 'ethereum/migrations',
  // contracts_build_directory: 'ethereum/build/contracts',

  networks: {
    ganache: {
      host: "127.0.0.1",
      port: 8545,
      network_id: "*"
    }
  }

};
