import * as ethers from 'ethers'

const truffle = new ethers.providers.JsonRpcProvider('http://localhost:9545')
const contract = require('../../build/contracts/TTT.json')
const ttt = new ethers.Contract('0x8f0483125fcb9aaaefa9209d8e9d7b9c8b9fb90f', contract.abi, truffle)

console.log(ttt)

const mState = {
  nonce: 0,
  tag: 0,
  data: '0x',
  state: {
    tag: 0,
    data: '0x000000000000000000'
  }
}

const mMove = {
  playerID: 0,
  data: '0x04'
}

console.log(mState)

ttt.nextState1(mState, mMove).then((mState) => {
  console.log(mState[0])

  const mMove = {
    playerID: 1,
    data: '0x00'
  }

  return ttt.nextState1(mState[0], mMove)

}).then((mState) => {
  console.log(mState[0])

})

const game = '0x0123456789012345678901234567890123456789'
const matchID = 12345
const matchIDHex = ethers.utils.hexlify(ethers.utils.padZeros(ethers.utils.arrayify(ethers.utils.hexlify(matchID)), 4))
const subkey = ethers.Wallet.createRandom()
const message = `Sign to play! This won't cost anything.\n\nGame: ${game}\nMatch: ${matchIDHex}\nPlayer: ${subkey.getAddress().toLowerCase()}`
const metamask = new ethers.providers.Web3Provider((global as any).web3.currentProvider)

metamask.listAccounts().then((accounts) => {
  metamask.getSigner(accounts[0]).signMessage(message).then((signature) => {
    const signatureBytes = ethers.utils.arrayify(signature)

    const dMatch = {
      game: game,
      matchID: matchID,
      players: [
        {
          account: accounts[0],
          subkey: subkey.getAddress(),
          subkeySignature: {
            v: signatureBytes[64],
            r: ethers.utils.hexlify(new Uint8Array(signatureBytes.buffer, 0, 32)),
            s: ethers.utils.hexlify(new Uint8Array(signatureBytes.buffer, 32, 32))
          }
        },
        {
          account: accounts[0],
          subkey: subkey.getAddress(),
          subkeySignature: {
            v: signatureBytes[64],
            r: ethers.utils.hexlify(new Uint8Array(signatureBytes.buffer, 0, 32)),
            s: ethers.utils.hexlify(new Uint8Array(signatureBytes.buffer, 32, 32))
          }
        }
      ],
      signature: {
        v: 0,
        r: '0x0000000000000000000000000000000000000000000000000000000000000000',
        s: '0x0000000000000000000000000000000000000000000000000000000000000000'
      }
    }

    dMatch['[object Object]'] = dMatch.players // XXX

    ttt.isSubkeySigned(dMatch, 0).then((signed) => {
      console.log(`ttt.isSubkeySigned: ${signed}`)
    })
  })
})
