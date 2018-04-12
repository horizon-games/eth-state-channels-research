import * as config from './config.js'

const cfg = config[process.env.NODE_ENV]

const arcadeumAddress = cfg.arcadeumAddress
const gameAddress = cfg.gameAddress
const serverAddress = cfg.serverAddress
const deposit = cfg.deposit

export { arcadeumAddress, gameAddress, serverAddress, deposit }
