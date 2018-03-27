import * as ethers from 'ethers'

const provider = new ethers.providers.JsonRpcProvider('http://localhost:8545')
const privateKey1 = '0x829e924fdf021ba3dbbc4225edfece9aca04b929d6e75613329ca6f1d31c0bb4'
const privateKey2 = '0xb0057716d5917badaf911b193b12b910811c1497b5bada8d7711f758981c3773'
const wallet1 = new ethers.Wallet(privateKey1, provider)
const wallet2 =  new ethers.Wallet(privateKey2, provider)

export { wallet1, wallet2 }
