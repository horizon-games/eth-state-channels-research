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

ttt.nextState(mState, mMove).then((mState) => {
  console.log(mState)

  mState = {
    nonce: mState[0].nonce,
    tag: mState[0].tag,
    data: mState[0].data,
    state: {
      tag: mState[0].state.tag,
      data: mState[0].state.data
    }
  }

  const mMove = {
    playerID: 1,
    data: '0x00'
  }

  return ttt.nextState(mState, mMove)

}).then((mState) => {
  console.log(mState)

})
