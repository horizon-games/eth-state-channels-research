import * as wsrelay from '../wsrelay'
import * as ethers from 'ethers'
import * as rxjs from 'rxjs'

const ArcadeumAddress = `0xcfeb869f69431e42cdb54a4f4f105c19c080a601`
const ServerAddress = `ws://localhost:8000/`

const ArcadeumContract = require(`arcadeum-contracts/build/contracts/Arcadeum.json`)
const GameContract = require(`arcadeum-contracts/build/contracts/DGame.json`)

export class Game {
  constructor(gameAddress: string, options?: { arcadeumAddress?: string, serverAddress?: string, wallet?: ethers.Wallet }) {
    let arcadeumAddress = ArcadeumAddress
    if (options !== undefined && options.arcadeumAddress !== undefined) {
      arcadeumAddress = options.arcadeumAddress
    }

    this.serverAddress = ServerAddress
    if (options !== undefined && options.serverAddress !== undefined) {
      this.serverAddress = options.serverAddress
    }

    if (options !== undefined && options.wallet !== undefined) {
      this.signer = options.wallet

    } else {
      const web3 = (window as any)[`web3`]
      const provider = new ethers.providers.Web3Provider(web3.currentProvider)
      this.signer = provider.getSigner()
    }

    this.arcadeumContract = new ethers.Contract(arcadeumAddress, ArcadeumContract.abi, this.signer)
    this.gameContract = new ethers.Contract(gameAddress, GameContract.abi, this.signer)
  }

  async deposit(wei: ethers.utils.BigNumber): Promise<string> {
    return (await this.arcadeumContract.deposit({ value: wei })).hash
  }

  createMatch(secretSeed: Uint8Array): Match {
    return new BasicMatch(secretSeed, this.arcadeumContract, this.gameContract, this.serverAddress, this.signer)
  }

  private readonly arcadeumContract: ethers.Contract
  private readonly gameContract: ethers.Contract
  private readonly serverAddress: string
  private readonly signer: ethers.Wallet | ethers.providers.Web3Signer
}

export interface Match {
  readonly ready: Promise<void>
  readonly playerID?: number
  readonly opponentID?: number
  readonly state: Promise<State>

  addCallback(callback: NextStateCallback): void
  createMove(data: Uint8Array): Promise<Move>
  queueMove(move: Move): void
}

export interface State {
  readonly winner: Promise<Winner>
  readonly nextPlayers: Promise<NextPlayers>

  isMoveLegal(move: Move): Promise<{ isMoveLegal: boolean, reason: number }>
  nextState(aMove: Move | [Move] | [Move, Move], anotherMove?: Move): Promise<State>
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

export interface Move {
  readonly playerID: number
  readonly data: Uint8Array
}

export interface NextStateCallback {
  (nextState: State, previousState?: State, aMove?: Move, anotherMove?: Move): Promise<void>
}

class BasicMatch implements Match, rxjs.Observer<wsrelay.Message> {
  constructor(private readonly secretSeed: Uint8Array, private readonly arcadeumContract: ethers.Contract, private readonly gameContract: ethers.Contract, private readonly serverAddress: string, private readonly signer: ethers.Wallet | ethers.providers.Web3Signer) {
    this.callbacks = []
    this.queue = []
    this.isRunning = true
    this.didQueueChange = false
    this.processedMoves = [undefined, undefined]
    this.playerMoves = []
    this.next = (message: wsrelay.Message) => { this.queueMove(new BasicMove(JSON.parse(message.payload))) }
    this.error = (error: any) => {}
    this.complete = () => {}
  }

  get ready(): Promise<void> {
    if (this.readyPromise === undefined) {
      this.readyPromise = this.getReady()
    }

    return this.readyPromise
  }

  playerID?: number

  get opponentID(): number | undefined {
    if (this.playerID === undefined) {
      return undefined
    }

    return 1 - this.playerID
  }

  get state(): Promise<BasicState> {
    return this.ready.then(() => this.statePromise!)
  }

  addCallback(callback: NextStateCallback): void {
    this.callbacks.push(callback)
  }

  async createMove(data: Uint8Array): Promise<Move> {
    if (this.playerID === undefined || this.subkey === undefined) {
      throw Error(`match not ready`)
    }

    const move = new BasicMove({ playerID: this.playerID, data: data })
    await move.sign(this.subkey, await this.state)
    return move
  }

  queueMove(move: Move): void {
    this.queue.push(move)

    if (this.isRunning) {
      this.didQueueChange = true
    } else {
      this.processQueue()
    }
  }

  readonly next: (message: wsrelay.Message) => void
  readonly error: (error: any) => void
  readonly complete: () => void

  private game?: string
  private timestamp?: ethers.utils.BigNumber
  private players?: [Player, Player]
  private matchSignature?: Signature
  private opponentSubkeySignature?: Signature

  private get opponentTimestampSignature(): Signature | undefined {
    if (this.opponentID === undefined || this.players === undefined) {
      return undefined
    }

    return this.players[this.opponentID].timestampSignature
  }

  private readyPromise?: Promise<void>
  private statePromise?: Promise<BasicState>
  private relay?: wsrelay.Relay
  private subkey?: ethers.Wallet
  private callbacks: NextStateCallback[]
  private queue: Move[]
  private isRunning: boolean
  private didQueueChange: boolean
  private processedMoves: [Move | undefined, Move | undefined]
  private agreedState?: BasicState
  private opponentMove?: Move
  private playerMoves: Move[]
  private random?: Uint8Array

  private async getReady(): Promise<void> {
    this.subkey = ethers.Wallet.createRandom()
    const subkeyMessage = await this.arcadeumContract.subkeyMessage(this.subkey.address)
    const subkeySignature = new Signature(await this.signer.signMessage(subkeyMessage))

    const relaySignature = new wsrelay.Signature(subkeySignature.v, base64(subkeySignature.r), base64(subkeySignature.s))
    this.relay = new wsrelay.Relay(this.serverAddress, base64(this.secretSeed), relaySignature, this.subkey.address, 1)
    this.relay.subscribe(this)

    const timestamp = JSON.parse((await this.relay.connectForTimestamp()).payload)
    const timestampSignature = sign(this.subkey, [`uint`], [timestamp])

    this.relay.send(JSON.stringify({
      gameID: 1,
      subkey: this.subkey.address,
      timestamp: timestamp,
      signature: timestampSignature
    }), 2)

    const response = JSON.parse((await this.relay.connectForMatchVerified()).payload)
    response.players[0].publicSeed = unbase64(response.players[0].publicSeed)
    response.players[1].publicSeed = unbase64(response.players[1].publicSeed)
    response.players[0].timestampSignature = new Signature(response.players[0].timestampSignature)
    response.players[1].timestampSignature = new Signature(response.players[1].timestampSignature)
    response.matchSignature = new Signature(response.matchSignature)
    response.opponentSubkeySignature = new Signature(response.opponentSubkeySignature)

    this.game = response.game
    this.timestamp = response.timestamp
    this.playerID = response.playerID
    this.players = response.players
    this.matchSignature = response.matchSignature
    this.opponentSubkeySignature = response.opponentSubkeySignature

    const initialState = this.gameContract.initialState(this.players![0].publicSeed, this.players![1].publicSeed)
    this.statePromise = initialState.then((response: MetaState) => {
      const state = new BasicState(response, this.arcadeumContract, this.gameContract)
      this.agreedState = state
      return state
    })

    this.statePromise!.then((state: BasicState) => this.runCallbacks(state))
    this.isRunning = true

    // XXX
    // @ts-ignore
    this[`[object Object]`] = this.players
  }

  private async runCallbacks(nextState: BasicState, previousState?: BasicState, aMove?: Move, anotherMove?: Move): Promise<void> {
    this.isRunning = true

    switch (nextState.metaState.tag) {
    case MetaTag.CommittingRandom:
      this.random = ethers.utils.randomBytes(nextState.metaState.data[0])
      this.queueMove(await this.createMove(ethers.utils.arrayify(ethers.utils.keccak256(this.random))))
      break

    case MetaTag.RevealingRandom:
      this.queueMove(await this.createMove(this.random!))
      delete this.random
      break

    default:
      const run = (callback: NextStateCallback) => callback(nextState, previousState, aMove, anotherMove).catch((reason: any) => {})
      await Promise.all(this.callbacks.map(run))
      break
    }

    this.processQueue()
  }

  private async processQueue(): Promise<void> {
    this.isRunning = true

    const [
      { state: state, stateHash: stateHash },
      opponent
    ] = await Promise.all([
      this.state.then(async (state: BasicState) => ({ state: state, stateHash: await state.hash })),
      this.arcadeumContract.playerAccount(this.timestamp, this.opponentTimestampSignature, this.opponentSubkeySignature)
    ])

    const queue = [...this.queue]
    this.didQueueChange = false
    const canProcessMove = await Promise.all(queue.map((move: Move) => this.canProcessMove(move, state, stateHash, opponent)))
    const movesToProcess = queue.filter((move: Move, i: number) => canProcessMove[i])
    this.queue = this.queue.filter((move: Move) => movesToProcess.indexOf(move) === -1)

    for (let move of movesToProcess) {
      await this.processMove(move as BasicMove)
    }

    if (await this.state === state) {
      if (this.didQueueChange) {
        this.processQueue()
      } else {
        this.isRunning = false
      }
    }
  }

  private async canProcessMove(move: Move, state: BasicState, stateHash: Uint8Array, opponent: string): Promise<boolean> {
    if (move instanceof BasicMove) {
      if (move.playerID === this.playerID) {
        if (move.hasStateHash(stateHash)) {
          return true
        }
      }

      if (move.playerID === this.opponentID) {
        const moveMaker = await this.arcadeumContract.moveMaker(state.metaState, move, this.opponentSubkeySignature)

        if (moveMaker === opponent) {
          return true
        }
      }
    }

    return false
  }

  private async processMove(move: BasicMove): Promise<void> {
    const state = await this.state

    if (move.playerID === this.playerID) {
      const stateHash = await state.hash

      if (!move.hasStateHash(stateHash)) {
        throw Error(`move not signed by player`)
      }

      this.relay!.send(JSON.stringify(move))

    } else /* move.playerID === this.opponentID */ {
      const [
        opponent,
        moveMaker
      ] = await Promise.all([
        this.arcadeumContract.playerAccount(this.timestamp, this.opponentTimestampSignature, this.opponentSubkeySignature),
        this.arcadeumContract.moveMaker(state.metaState, move, this.opponentSubkeySignature)
      ])

      if (moveMaker !== opponent) {
        throw Error(`move not signed by opponent`)
      }

      const { isMoveLegal: isMoveLegal, reason: reason } = await state.isMoveLegal(move)

      if (!isMoveLegal) {
        if (await this.arcadeumContract.canReportCheater(this, state.metaState, move)) {
          this.arcadeumContract.reportCheater(this, state.metaState, move)
        }

        throw Error(`illegal opponent move: reason ${reason}`)
      }
    }

    if (this.processedMoves[move.playerID] !== undefined) {
      throw Error(`already processed player ${move.playerID}'s move`)
    }

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
      this.runCallbacks(await this.state, state, move)
      break

    case NextPlayers.Both:
      this.processedMoves[move.playerID] = move

      if (this.processedMoves[0] === undefined || this.processedMoves[1] === undefined) {
        return
      }

      const processedMoves = this.processedMoves as [Move, Move]
      this.processedMoves = [undefined, undefined]
      this.agreedState = state
      this.opponentMove = processedMoves[this.opponentID!]
      this.playerMoves = [processedMoves[this.playerID!]]
      this.statePromise = state.nextState(processedMoves)
      this.runCallbacks(await this.state, state, processedMoves[0], processedMoves[1])
      break
    }

    const winner = await (await this.state).winner

    if (winner === Winner.Player0 && this.playerID === 0 || winner === Winner.Player1 && this.playerID === 1) {
      if (await this.arcadeumContract.canClaimReward(this, this.agreedState!.metaState, this.opponentMove, this.playerMoves)) {
        this.arcadeumContract.claimReward(this, this.agreedState!.metaState, this.opponentMove, this.playerMoves)
      }
    }
  }
}

interface Player {
  readonly seedRating: number
  readonly publicSeed: Uint8Array
  readonly timestampSignature: Signature
}

interface MetaState {
  readonly nonce: number
  readonly tag: MetaTag
  readonly data: Uint8Array
  readonly state: {
    readonly tag: number
    readonly data: Uint8Array
  }
}

enum MetaTag {
  None,
  CommittingRandom,
  RevealingRandom,
  CommittingSecret,
  RevealingSecret
}

class BasicState implements State {
  constructor(metaState: MetaState, private readonly arcadeumContract: ethers.Contract, private readonly gameContract: ethers.Contract) {
    this.tag = metaState.state.tag
    this.data = ethers.utils.arrayify(metaState.state.data)
    this.meta = {
      nonce: metaState.nonce,
      tag: metaState.tag,
      data: ethers.utils.arrayify(metaState.data)
    }
  }

  get winner(): Promise<Winner> {
    return this.gameContract.winner(this.metaState)
  }

  get nextPlayers(): Promise<NextPlayers> {
    return this.gameContract.nextPlayers(this.metaState)
  }

  async isMoveLegal(move: Move): Promise<{ isMoveLegal: boolean, reason: number }> {
    const response = await this.gameContract.isMoveLegal(this.metaState, move)
    return { isMoveLegal: response[0], reason: response[1] }
  }

  async nextState(aMove: Move | [Move] | [Move, Move], anotherMove?: Move): Promise<BasicState> {
    let nextState: (metaState: MetaState, aMove: Move, anotherMove?: Move) => Promise<MetaState>

    if (aMove instanceof Array) {
      if (anotherMove !== undefined) {
        throw Error(`unexpected second argument: array already given`)
      }

      switch (aMove.length) {
      case 1:
        nextState = this.gameContract[`nextState((uint32,uint8,bytes,(uint32,bytes)),(uint8,bytes))`]
        return new BasicState(await nextState(this.metaState, aMove[0]), this.arcadeumContract, this.gameContract)

      case 2:
        nextState = this.gameContract[`nextState((uint32,uint8,bytes,(uint32,bytes)),(uint8,bytes),(uint8,bytes))`]
        return new BasicState(await nextState(this.metaState, aMove[0], aMove[1]), this.arcadeumContract, this.gameContract)
      }

    } else /* aMove: Move */ {
      if (anotherMove === undefined) {
        nextState = this.gameContract[`nextState((uint32,uint8,bytes,(uint32,bytes)),(uint8,bytes))`]
        return new BasicState(await nextState(this.metaState, aMove), this.arcadeumContract, this.gameContract)

      } else {
        nextState = this.gameContract[`nextState((uint32,uint8,bytes,(uint32,bytes)),(uint8,bytes),(uint8,bytes))`]
        return new BasicState(await nextState(this.metaState, aMove, anotherMove), this.arcadeumContract, this.gameContract)
      }
    }

    throw Error(`expected dgame.Move[] of length 1 or 2`)
  }

  get hash(): Promise<Uint8Array> {
    const hashPromise = this.arcadeumContract.stateHash(this.metaState)
    return hashPromise.then((response: string) => ethers.utils.arrayify(response))
  }

  get metaState(): MetaState {
    return {
      nonce: this.meta.nonce,
      tag: this.meta.tag,
      data: this.meta.data,
      state: {
        tag: this.tag,
        data: this.data
      }
    }
  }

  private readonly tag: number
  private readonly data: Uint8Array
  private readonly meta: {
    readonly nonce: number
    readonly tag: MetaTag
    readonly data: Uint8Array
  }
}

class BasicMove implements Move {
  constructor(readonly move: { playerID: number, data: string | Uint8Array, signature?: Signature }) {
    if (typeof move.data === `string`) {
      move.data = unbase64(move.data)
    }

    this.playerID = move.playerID
    this.data = move.data
    this.signature = new Signature(move.signature)
  }

  readonly playerID: number
  readonly data: Uint8Array

  async sign(subkey: ethers.Wallet, state: BasicState): Promise<void> {
    const { isMoveLegal: isMoveLegal, reason: reason } = await state.isMoveLegal(this)

    if (!isMoveLegal) {
      throw Error(`illegal player move: reason ${reason}`)
    }

    this.stateHash = await state.hash
    const types = [`bytes32`, `uint8`, `bytes`]
    const values = [this.stateHash, this.playerID, this.data]
    this.signature = sign(subkey, types, values)
  }

  hasStateHash(stateHash: Uint8Array): boolean {
    return this.stateHash !== undefined && areArraysEqual(stateHash, this.stateHash)
  }

  private stateHash?: Uint8Array
  private signature?: Signature

  private toJSON(): any {
    return {
      playerID: this.playerID,
      data: base64(this.data),
      signature: this.signature
    }
  }
}

class Signature {
  constructor(signature?: string | { readonly v: number, readonly r: string | Uint8Array, readonly s: string | Uint8Array }) {
    if (typeof signature === `string`) {
      const signatureBytes = ethers.utils.arrayify(signature)

      this.v = signatureBytes[64]
      this.r = new Uint8Array(signatureBytes.buffer, 0, 32)
      this.s = new Uint8Array(signatureBytes.buffer, 32, 32)

    } else if (signature !== undefined) {
      this.v = signature.v

      if (typeof signature.r === `string`) {
        this.r = unbase64(signature.r)
      } else {
        this.r = signature.r
      }

      if (typeof signature.s === `string`) {
        this.s = unbase64(signature.s)
      } else {
        this.s = signature.s
      }

    } else {
      this.v = 0
      this.r = new Uint8Array(32)
      this.s = new Uint8Array(32)
    }
  }

  readonly v: number
  readonly r: Uint8Array
  readonly s: Uint8Array

  private toJSON(): any {
    return {
      v: this.v,
      r: base64(this.r),
      s: base64(this.s)
    }
  }
}

function sign(wallet: ethers.Wallet, types: string[], values: any[]): Signature {
  const hash = ethers.utils.solidityKeccak256(types, values)
  const signature = new ethers.SigningKey(wallet.privateKey).signDigest(hash)

  return new Signature({
    v: 27 + signature.recoveryParam,
    r: ethers.utils.padZeros(ethers.utils.arrayify(signature.r), 32),
    s: ethers.utils.padZeros(ethers.utils.arrayify(signature.s), 32)
  })
}

function base64(data: Uint8Array): string {
  // XXX
  // @ts-ignore
  return Buffer.from(data).toString(`base64`)
}

function unbase64(data: string): Uint8Array {
  return new Uint8Array(Buffer.from(data, `base64`))
}

interface Arrayish {
  readonly length: number
  readonly [index: number]: any
}

function areArraysEqual(anArray: Arrayish, anotherArray: Arrayish): boolean {
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
