import * as ethers from 'ethers'

const provider = new ethers.providers.JsonRpcProvider('http://localhost:8545')
const privateKey1 = '0x4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d'
const privateKey2 = '0x6cbed15c793ce57650b9877cf6fa156fbef513c4e6134f022a85b1ffdd59b2a1'
const wallet1 = new ethers.Wallet(privateKey1, provider)
const wallet2 =  new ethers.Wallet(privateKey2, provider)

export { wallet1, wallet2 }
