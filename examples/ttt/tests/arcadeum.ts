import * as config from './config.js'

const cfg = config[process.env.NODE_ENV]

const arcadeumAddress = cfg.arcadeumAddress
const gameAddress = cfg.gameAddress
const arcadeumServerHost = cfg.arcadeumServerHost
const arcadeumServerPort = cfg.arcadeumServerPort
const deposit = cfg.deposit

export { arcadeumAddress, gameAddress, arcadeumServerHost, arcadeumServerPort, deposit }