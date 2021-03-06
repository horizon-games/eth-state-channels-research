module.exports = {
  development: {
    network: "ganache",
    arcadeumAddress: "0xcfeb869f69431e42cdb54a4f4f105c19c080a601",
    gameAddress: "0x9561c133dd8580860b6b7e504bc5aa500f0f06a7",
    wallet1Password: "0x829e924fdf021ba3dbbc4225edfece9aca04b929d6e75613329ca6f1d31c0bb4",
    wallet2Password: "0xb0057716d5917badaf911b193b12b910811c1497b5bada8d7711f758981c3773",
    deposit: "0.02",
    serverAddress: "ws://localhost:8000/",
    jsonRpcUrl: "http://localhost:8545"
  },
  staging: {
    network: "rinkeby",
    arcadeumAddress: "0x29de34e0f36813f140c80788dbb0faeae38fdd94",
    gameAddress: "0xa61cb9020b81f721594d9dcba3abc8d39395d7f2",
    wallet1Password: "0x4a05e4e2c8b80906ccf688008b8f257fa984f05c27260bae0d78515f38d6f412",
    wallet2Password: "0x5D862464FE9303452126C8BC94274B8C5F9874CBD219789B3EB2128075A76F72",
    deposit: "0.02",
    serverAddress: "wss://relay.arcadeum.com:80/",
    infuraApiToken: "P8djn1ELvrq7uw7LrE22"
  }
}
