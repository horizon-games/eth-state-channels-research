import * as ethers from 'ethers'

const provider = new ethers.providers.JsonRpcProvider('http://localhost:9545')
const contract = require('../../build/contracts/TTT.json')
const ttt = new ethers.Contract('0x8f0483125fcb9aaaefa9209d8e9d7b9c8b9fb90f', contract.abi, provider)

console.log(ttt)

const state = {
  nonce: 0,
  type_: 0,
  data: new Uint8Array(0),
  statuses: [],
  state: {
    type_: 0,
    data: new Uint8Array(0),
    statuses: []
  }
}

ttt.playerStatus(state, 0).then((status_) => {
  console.log(`Player 1's status is ${status_}`)
})

ttt.playerStatus(state, 1).then((status_) => {
  console.log(`Player 2's status is ${status_}`)
})
