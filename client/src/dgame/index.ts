import * as ethers from 'ethers'
import * as wsrelay from '../wsrelay'

export class DGame {
  constructor(arcadeumAddress: string, gameAddress: string, options: { arcadeumServerHost?: string, arcadeumServerPort?: number, account?: ethers.Wallet, ssl: boolean } = { ssl: false }) {
    this.arcadeumServerHost = options.arcadeumServerHost !== undefined ? options.arcadeumServerHost : 'localhost'
    this.arcadeumServerPort = options.arcadeumServerPort !== undefined ? options.arcadeumServerPort : 8000
    this.account = options.account
    this.ssl = options.ssl

    const arcadeumMetadata = require('arcadeum-contracts/build/contracts/Arcadeum.json')
    const gameMetadata = require('arcadeum-contracts/build/contracts/DGame.json')

    if (this.account !== undefined) {
      this.arcadeumContract = new ethers.Contract(arcadeumAddress, arcadeumMetadata.abi, this.account)
      this.gameContract = new ethers.Contract(gameAddress, gameMetadata.abi, this.account)

    } else {
      this.signer = (new ethers.providers.Web3Provider((window as any).web3.currentProvider)).getSigner() // XXX: choose account
      this.arcadeumContract = new ethers.Contract(arcadeumAddress, arcadeumMetadata.abi, this.signer)
      this.gameContract = new ethers.Contract(gameAddress, gameMetadata.abi, this.signer)
    }
  }

  get address(): string {
    return this.gameContract.address
  }

  async deposit(value: ethers.utils.BigNumber): Promise<{ hash: string }> {
    return this.arcadeumContract.deposit({ value: value })
  }

  get matchDuration(): Promise<number> {
    return this.gameContract.matchDuration().then((response: ethers.utils.BigNumber) => response.toNumber())
  }

  async isSecretSeedValid(address: string, secretSeed: Uint8Array): Promise<boolean> {
    return this.gameContract.isSecretSeedValid(address, secretSeed)
  }

  async createMatch(secretSeed: Uint8Array, onNextState?: NextStateCallback): Promise<Match> {
    const subkey = ethers.Wallet.createRandom()
    const subkeyMessage = await this.arcadeumContract.subkeyMessage(subkey.getAddress())

    let subkeySignature: Signature
    if (this.account !== undefined) {
      subkeySignature = new Signature(await this.account.signMessage(subkeyMessage))
    } else /* this.signer !== undefined */ {
      subkeySignature = new Signature(await this.signer!.signMessage(subkeyMessage))
    }

    const seed64 = base64(secretSeed)
    const r64 = base64(subkeySignature.r)
    const s64 = base64(subkeySignature.s)
    const relay = new wsrelay.Relay(this.arcadeumServerHost, this.arcadeumServerPort, this.ssl, seed64, new wsrelay.Signature(subkeySignature.v, r64, s64), subkey.getAddress(), 1)
    const timestamp = JSON.parse((await relay.connectForTimestamp()).payload)
    const timestampSignature = sign(subkey, [`uint`], [timestamp])

    relay.send(JSON.stringify({
      gameID: 1,
      subkey: subkey.getAddress(),
      timestamp: timestamp,
      signature: {
        v: timestampSignature.v,
        r: base64(timestampSignature.r),
        s: base64(timestampSignature.s)
      }
    }), 2)

    const response = JSON.parse((await relay.connectForMatchVerified()).payload)

    response.players[0].publicSeed = [ethers.utils.bigNumberify(unbase64(response.players[0].publicSeed))]
    response.players[1].publicSeed = [ethers.utils.bigNumberify(unbase64(response.players[1].publicSeed))]
    response.players[0].timestampSignature.r = ethers.utils.arrayify(ethers.utils.toUtf8String(unbase64(response.players[0].timestampSignature.r)))
    response.players[1].timestampSignature.r = ethers.utils.arrayify(ethers.utils.toUtf8String(unbase64(response.players[1].timestampSignature.r)))
    response.players[0].timestampSignature.s = ethers.utils.arrayify(ethers.utils.toUtf8String(unbase64(response.players[0].timestampSignature.s)))
    response.players[1].timestampSignature.s = ethers.utils.arrayify(ethers.utils.toUtf8String(unbase64(response.players[1].timestampSignature.s)))
    response.matchSignature.r = unbase64(response.matchSignature.r)
    response.matchSignature.s = unbase64(response.matchSignature.s)
    response.opponentSubkeySignature.r = ethers.utils.arrayify(ethers.utils.toUtf8String(unbase64(response.opponentSubkeySignature.r)))
    response.opponentSubkeySignature.s = ethers.utils.arrayify(ethers.utils.toUtf8String(unbase64(response.opponentSubkeySignature.s)))

    return new RemoteMatch(response, subkey, this.arcadeumContract, this.gameContract, relay, onNextState)
  }

  private signer?: ethers.providers.Web3Signer
  private arcadeumContract: ethers.Contract
  private gameContract: ethers.Contract
  private arcadeumServerHost: string
  private arcadeumServerPort: number
  private account: ethers.Wallet | undefined
  private ssl: boolean
}

export interface Match {
  readonly playerID: number
  readonly opponentID: number
  readonly state: Promise<State>
  createMove(data: Uint8Array): Promise<Move>
  queueMove(move: Move): void
}

export interface NextStateCallback {
  (match: Match, previousState: State, nextState: State, aMove: Move, anotherMove?: Move): void
}

export interface State {
  readonly winner: Promise<Winner>
  readonly nextPlayers: Promise<NextPlayers>
  isMoveLegal(move: Move): Promise<{ isMoveLegal: boolean, reason: number }>
  nextState(aMove: Move, anotherMove?: Move): Promise<State>
  nextState(moves: [Move] | [Move, Move]): Promise<State>
  readonly encoding: any
  readonly hash: Promise<Uint8Array>
}

export enum Winner {
  None,
  Player0,
  Player1
}

export enum NextPlayers {
  None,
  Player0,
  Player1,
  Both
}

export class Move {
  constructor(readonly move: { playerID: number, data: Uint8Array, signature?: any }) {
    this.playerID = move.playerID
    this.data = move.data

    if (move.signature !== undefined) {
      this.signature = move.signature
    } else {
      this.signature = new Signature()
    }
  }

  async sign(subkey: ethers.Wallet, state: State): Promise<void> {
    this.stateHash = await state.hash
    this.signature = sign(subkey, [`bytes32`, `uint8`, `bytes`], [this.stateHash, this.playerID, this.data])
  }

  async wasSignedWithState(state: State): Promise<boolean> {
    if (this.stateHash === undefined) {
      return false
    }

    return areArraysEqual(await state.hash, this.stateHash)
  }

  readonly playerID: number
  readonly data: Uint8Array
  private stateHash?: Uint8Array
  private signature?: any
}

interface MatchInterface {
  readonly game: string
  readonly timestamp: ethers.utils.BigNumber
  readonly playerID: number
  readonly players: [PlayerInterface, PlayerInterface]
  readonly matchSignature: Signature
  readonly opponentSubkeySignature: Signature
}

interface PlayerInterface {
  readonly seedRating: number
  // XXX: https://github.com/ethereum/solidity/issues/3270
  readonly publicSeed: [Uint8Array]
  readonly timestampSignature: Signature
}

interface StateInterface {
  readonly nonce: number
  readonly tag: number
  // XXX: https://github.com/ethereum/solidity/issues/3270
  readonly data: [Uint8Array, Uint8Array, Uint8Array]
  readonly state: {
    readonly tag: number
    // XXX: https://github.com/ethereum/solidity/issues/3270
    readonly data: [Uint8Array]
  }
}

class BasicMatch implements Match {
  constructor(match: MatchInterface, private subkey: ethers.Wallet, private arcadeumContract: ethers.Contract, gameContract: ethers.Contract, onNextState?: NextStateCallback) {
    if (onNextState !== undefined) {
      this.onNextState = onNextState
    } else {
      this.onNextState = (match: Match, previousState: State, nextState: State, aMove: Move, anotherMove?: Move) => {}
    }

    this.game = match.game
    this.timestamp = match.timestamp
    this.playerID = match.playerID
    this.players = match.players
    this.matchSignature = match.matchSignature
    this.opponentSubkeySignature = match.opponentSubkeySignature

    this.statePromise = gameContract.initialState(this.players[0].publicSeed, this.players[1].publicSeed).then((response: StateInterface) => {
      const state = new BasicState(response, arcadeumContract, gameContract)
      this.agreedState = state
      return state
    })

    this.playerMoves = []
    this.appliedMoves = [undefined, undefined]
    this.queue = []
    this.isProcessingQueue = false

    // @ts-ignore
    this[`[object Object]`] = this.players // XXX
  }

  readonly playerID: number

  get opponentID(): number {
    return 1 - this.playerID
  }

  get state(): Promise<BasicState> {
    return this.getState()
  }

  async createMove(data: Uint8Array): Promise<Move> {
    const move = new Move({ playerID: this.playerID, data: data })
    const state = await this.statePromise

    const {
      isMoveLegal: isMoveLegal,
      reason: reason
    } = await state.isMoveLegal(move)

    if (!isMoveLegal) {
      throw Error(`illegal player move: reason ${reason}`)
    }

    await move.sign(this.subkey, state)
    return move
  }

  queueMove(move: Move): void {
    if (this.queue.indexOf(move) !== -1) {
      return
    }

    this.queue.push(move)
    this.processQueue()
  }

  protected getState(): Promise<BasicState> {
    return this.statePromise
  }

  protected onNextState: NextStateCallback

  protected async applyMove(move: Move): Promise<boolean> {
    const state = await this.statePromise

    if (move.playerID === this.playerID) {
      if (!await move.wasSignedWithState(state)) {
        return false
      }

    } else /* move.playerID === this.opponentID */ {
      const [
        opponent,
        moveMaker
      ] = await Promise.all([
        this.arcadeumContract.playerAccount(this.timestamp, this.opponentTimestampSignature, this.opponentSubkeySignature),
        this.arcadeumContract.moveMaker(state.encoding, move, this.opponentSubkeySignature)
      ])

      if (moveMaker !== opponent) {
        return false
      }

      const {
        isMoveLegal: isMoveLegal,
        reason: reason
      } = await state.isMoveLegal(move)

      if (!isMoveLegal) {
        if (await this.arcadeumContract.canReportCheater(this, state.encoding, move)) {
          this.arcadeumContract.reportCheater(this, state.encoding, move)
        }

        throw Error(`illegal opponent move: reason ${reason}`)
      }
    }

    if (this.appliedMoves[move.playerID] !== undefined) {
      return false
    }

    let nextState: BasicState

    switch (await state.nextPlayers) {
    case NextPlayers.Player0:
    case NextPlayers.Player1:
      if (move.playerID === this.playerID) {
        this.playerMoves.push(move)

      } else /* move.playerID === this.opponentID */ {
        this.agreedState = state
        this.opponentMove = move
        this.playerMoves = []
      }

      this.statePromise = state.nextState(move)
      nextState = await this.statePromise
      this.onNextState(this, state, nextState, move)
      break

    case NextPlayers.Both:
      this.appliedMoves[move.playerID] = move

      if (this.appliedMoves[0] === undefined || this.appliedMoves[1] === undefined) {
        return true
      }

      const appliedMoves = [this.appliedMoves[0], this.appliedMoves[1]] as [Move, Move]
      this.appliedMoves = [undefined, undefined]
      this.agreedState = state
      this.opponentMove = appliedMoves[this.opponentID]
      this.playerMoves = [appliedMoves[this.playerID]]
      this.statePromise = state.nextState(appliedMoves)
      nextState = await this.statePromise
      this.onNextState(this, state, nextState, appliedMoves[0], appliedMoves[1])
      break

    default:
      throw Error(`impossible since move is legal`)
    }

    const winner = await nextState.winner

    if (winner === Winner.Player0 && this.playerID === 0 || winner === Winner.Player1 && this.playerID === 1) {
      if (await this.arcadeumContract.canClaimReward(this, this.agreedState!.encoding, this.opponentMove, this.playerMoves)) {
        this.arcadeumContract.claimReward(this, this.agreedState!.encoding, this.opponentMove, this.playerMoves)
      }
    }

    return true
  }

  private async processQueue(): Promise<void> {
    if (this.isProcessingQueue) {
      return
    }

    this.isProcessingQueue = true

    let queue: Move[]

    do {
      queue = [...this.queue]

      for (let move of queue) {
        if (await this.applyMove(move)) {
          this.queue = this.queue.filter((aMove: Move) => aMove !== move)
        }
      }
    } while (!areArraysEqual(this.queue, queue))

    this.isProcessingQueue = false
  }

  private readonly game: string
  private readonly timestamp: ethers.utils.BigNumber
  private readonly players: [PlayerInterface, PlayerInterface]
  private readonly matchSignature: Signature
  private readonly opponentSubkeySignature: Signature

  private get opponentTimestampSignature(): Signature {
    return this.players[this.opponentID].timestampSignature
  }

  private statePromise: Promise<BasicState>
  private agreedState?: BasicState
  private opponentMove?: Move
  private playerMoves: Move[]
  private appliedMoves: [Move | undefined, Move | undefined]
  private queue: Move[]
  private isProcessingQueue: boolean
}

class RemoteMatch extends BasicMatch {
  constructor(match: MatchInterface, subkey: ethers.Wallet, arcadeumContract: ethers.Contract, gameContract: ethers.Contract, private relay: wsrelay.Relay, onNextState?: NextStateCallback) {
    super(match, subkey, arcadeumContract, gameContract)

    super.getState().then(async (state: BasicState) => {
      switch (state.metadata.tag) {
      case MetaTag.CommittingRandom:
        if (this.random === undefined) {
          this.random = ethers.utils.randomBytes(state.metadata.data[0][31])
          this.queueMove(await this.createMove(ethers.utils.arrayify(ethers.utils.keccak256(this.random))))
        }

        break
      }
    })

    this.onNextState = async (match: Match, previousState: State, nextState: State, aMove: Move, anotherMove?: Move) => {
      if (!(match instanceof RemoteMatch)) {
        throw Error(`impossible`)
      }

      if (!(nextState instanceof BasicState)) {
        throw Error(`impossible`)
      }

      switch (nextState.metadata.tag) {
      case MetaTag.CommittingRandom:
        if (match.random === undefined) {
          match.random = ethers.utils.randomBytes(nextState.metadata.data[0][31])
          match.queueMove(await match.createMove(ethers.utils.arrayify(ethers.utils.keccak256(match.random))))
        }

        break

      case MetaTag.RevealingRandom:
        if (match.random !== undefined) {
          match.queueMove(await match.createMove(match.random))
          delete match.random
        }

        break

      default:
        if (onNextState !== undefined) {
          onNextState(match, previousState, nextState, aMove, anotherMove)
        }

        break
      }
    }

    relay.subscribe(this)
  }

  complete(): void {
  }

  error(error: any): void {
  }

  next(message: wsrelay.Message): void {
    const response = JSON.parse(message.payload)

    response.data = deserializeUint8Array(response.data)
    response.signature.r = deserializeUint8Array(response.signature.r)
    response.signature.s = deserializeUint8Array(response.signature.s)

    this.queueMove(new Move(response))
  }

  protected async applyMove(move: Move): Promise<boolean> {
    const wasApplied = await super.applyMove(move)

    if (wasApplied) {
      if (move.playerID === this.playerID) {
        this.relay.send(JSON.stringify(move))
      }
    }

    return wasApplied
  }

  private random?: Uint8Array
}

enum MetaTag {
  None,
  CommittingRandom,
  RevealingRandom,
  CommittingSecret,
  RevealingSecret
}

class BasicState implements State {
  constructor(state: StateInterface, private arcadeumContract: ethers.Contract, private gameContract: ethers.Contract) {
    this.tag = state.state.tag
    this.data = state.state.data
    this.metadata = {
      nonce: state.nonce,
      tag: state.tag,
      data: state.data
    }

    for (let i in this.data) {
      this.data[i] = ethers.utils.arrayify(this.data[i])
    }

    for (let i in this.metadata.data) {
      this.metadata.data[i] = ethers.utils.arrayify(this.metadata.data[i])
    }
  }

  readonly metadata: {
    readonly nonce: number
    readonly tag: number
    // XXX: https://github.com/ethereum/solidity/issues/3270
    readonly data: [Uint8Array, Uint8Array, Uint8Array]
  }

  get winner(): Promise<Winner> {
    return this.gameContract.winner(this.encoding)
  }

  get nextPlayers(): Promise<NextPlayers> {
    return this.gameContract.nextPlayers(this.encoding)
  }

  async isMoveLegal(move: Move): Promise<{ isMoveLegal: boolean, reason: number }> {
    const response = await this.gameContract.isMoveLegal(this.encoding, move)

    return {
      isMoveLegal: response[0],
      reason: response[1]
    }
  }

  async nextState(aMove: Move | [Move] | [Move, Move], anotherMove?: Move): Promise<BasicState> {
    if (aMove instanceof Array) {
      if (anotherMove !== undefined) {
        throw Error(`unexpected second argument: array already given`)
      }

      switch (aMove.length) {
      case 1:
        return new BasicState(await this.gameContract[`nextState((uint32,uint8,bytes32[3],(uint32,bytes32[1])),(uint8,bytes))`](this.encoding, aMove[0]), this.arcadeumContract, this.gameContract)

      case 2:
        return new BasicState(await this.gameContract[`nextState((uint32,uint8,bytes32[3],(uint32,bytes32[1])),(uint8,bytes),(uint8,bytes))`](this.encoding, aMove[0], aMove[1]), this.arcadeumContract, this.gameContract)
      }

    } else {
      if (anotherMove === undefined) {
        return new BasicState(await this.gameContract[`nextState((uint32,uint8,bytes32[3],(uint32,bytes32[1])),(uint8,bytes))`](this.encoding, aMove), this.arcadeumContract, this.gameContract)

      } else {
        return new BasicState(await this.gameContract[`nextState((uint32,uint8,bytes32[3],(uint32,bytes32[1])),(uint8,bytes),(uint8,bytes))`](this.encoding, aMove, anotherMove), this.arcadeumContract, this.gameContract)
      }
    }

    throw Error(`expected dgame.Move[] of length 1 or 2`)
  }

  get encoding(): any {
    return {
      nonce: this.metadata.nonce,
      tag: this.metadata.tag,
      data: this.metadata.data,
      state: {
        tag: this.tag,
        data: this.data
      }
    }
  }

  get hash(): Promise<Uint8Array> {
    return this.arcadeumContract.stateHash(this.encoding).then((response: string) => ethers.utils.arrayify(response))
  }

  private readonly tag: number
  // XXX: https://github.com/ethereum/solidity/issues/3270
  private readonly data: [Uint8Array]
}

class Signature {
  constructor(signature?: string | Signature) {
    if (typeof signature === `string`) {
      const signatureBytes = ethers.utils.arrayify(signature)

      this.v = signatureBytes[64]
      this.r = new Uint8Array(signatureBytes.buffer, 0, 32)
      this.s = new Uint8Array(signatureBytes.buffer, 32, 32)

    } else {
      if (signature !== undefined && signature.hasOwnProperty('v')) {
        this.v = signature.v
      } else {
        this.v = 0
      }

      if (signature !== undefined && signature.hasOwnProperty('r')) {
        this.r = signature.r
      } else {
        this.r = new Uint8Array(32)
      }

      if (signature !== undefined && signature.hasOwnProperty('s')) {
        this.s = signature.s
      } else {
        this.s = new Uint8Array(32)
      }
    }
  }

  readonly v: number
  readonly r: Uint8Array
  readonly s: Uint8Array
}

function sign(wallet: ethers.Wallet, types: string[], values: any[]): Signature {
  const hash = ethers.utils.solidityKeccak256(types, values)
  const signatureValues = new ethers.SigningKey(wallet.privateKey).signDigest(hash)

  return {
    v: 27 + signatureValues.recoveryParam,
    r: ethers.utils.padZeros(ethers.utils.arrayify(signatureValues.r), 32),
    s: ethers.utils.padZeros(ethers.utils.arrayify(signatureValues.s), 32)
  }
}

function deserializeUint8Array(data?: { [index: number]: number }): Uint8Array | undefined {
  if (data === undefined) {
    return undefined
  }

  const array: number[] = []

  for (let i = 0; data[i] !== undefined; i++) {
    array.push(data[i])
  }

  return new Uint8Array(array)
}

function base64(data: Uint8Array): string {
  return new Buffer(ethers.utils.hexlify(data)).toString(`base64`)
}

function unbase64(data: string): Uint8Array {
  return Uint8Array.from(Buffer.from(data, `base64`))
}

function areArraysEqual(anArray: { readonly length: number, [index: number]: any }, anotherArray: { readonly length: number, [index: number]: any }): boolean {
  if (anArray.length !== anotherArray.length) {
    return false
  }

  for (let i in anArray) {
    if (anArray[i] !== anotherArray[i]) {
      return false
    }
  }

  return true
}
