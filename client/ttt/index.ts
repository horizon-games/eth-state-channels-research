import * as ethers from 'ethers'

const provider = new ethers.providers.JsonRpcProvider('http://localhost:9545')
const contract = require('../../build/contracts/TTT.json')
const ttt = new ethers.Contract('0x8f0483125fcb9aaaefa9209d8e9d7b9c8b9fb90f', contract.abi, provider)

console.log(ttt)
