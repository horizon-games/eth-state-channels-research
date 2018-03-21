import * as ethers from 'ethers'
import * as wsrelay from 'wsrelay'

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

export class DGame {
  constructor(gameAddress: string, private account?: ethers.Wallet) {
    const arcadeumAddress = `0xcfeb869f69431e42cdb54a4f4f105c19c080a601`
    const arcadeumMetadata = require(`../../build/contracts/Arcadeum.json`)
    const gameMetadata = require(`../../build/contracts/DGame.json`)

    if (account !== undefined) {
      this.arcadeumContract = new ethers.Contract(arcadeumAddress, arcadeumMetadata.abi, account)
      this.gameContract = new ethers.Contract(gameAddress, gameMetadata.abi, account)

    } else {
      this.signer = (new ethers.providers.Web3Provider((window as any).web3.currentProvider)).getSigner() // XXX: choose account
      this.arcadeumContract = new ethers.Contract(arcadeumAddress, arcadeumMetadata.abi, this.signer)
      this.gameContract = new ethers.Contract(gameAddress, gameMetadata.abi, this.signer)
    }
  }

  get address(): string {
    return this.gameContract.address
  }

  async deposit(value: ethers.utils.BigNumber): Promise<void> {
    return this.arcadeumContract.deposit({ value: value })
  }

  get matchDuration(): Promise<number> {
    return this.gameContract.matchDuration().then(response => response.toNumber())
  }

  async isSecretSeedValid(address: string, secretSeed: Uint8Array): Promise<boolean> {
    return this.gameContract.isSecretSeedValid(address, secretSeed)
  }

  async createMatch(secretSeed: Uint8Array, onChange?: ChangeCallback, onCommit?: CommitCallback): Promise<Match> {
    const subkey = ethers.Wallet.createRandom()
    const subkeyMessage = await this.arcadeumContract.subkeyMessage(subkey.getAddress())

    let subkeySignature: Signature
    if (this.account !== undefined) {
      subkeySignature = sign(this.account, [`string`], subkeyMessage)
    } else /* this.signer !== undefined */ {
      subkeySignature = new Signature(await this.signer!.signMessage(subkeyMessage))
    }

    const seed64 = base64(secretSeed)
    const r64 = base64(subkeySignature.r)
    const s64 = base64(subkeySignature.s)
    const relay = new wsrelay.Relay(`localhost`, 8000, false, seed64, new wsrelay.Signature(subkeySignature.v, r64, s64), subkey.getAddress(), 1)
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

    return new RemoteMatch(relay, this.arcadeumContract, this.gameContract, subkey, response, onChange, onCommit)
  }

  private signer?: ethers.providers.Web3Signer
  private arcadeumContract: ethers.Contract
  private gameContract: ethers.Contract
}

export class Match {
  constructor(private arcadeumContract: ethers.Contract, private gameContract: ethers.Contract, private subkey: ethers.Wallet, match: MatchInterface, public onChange?: ChangeCallback, public onCommit?: CommitCallback) {
    this.game = match.game
    this.timestamp = match.timestamp
    this.playerID = match.playerID
    this.players = match.players
    this.matchSignature = match.matchSignature
    this.opponentSubkeySignature = match.opponentSubkeySignature
    this.playerMoves = []
    this.pendingMoves = [undefined, undefined]
    this[`[object Object]`] = this.players // XXX
  }

  get state(): Promise<State> {
    if (this.currentState === undefined) {
      return this.gameContract.initialState(this.players[0].publicSeed, this.players[1].publicSeed).then(response => {
        this.agreedState = new State(this.arcadeumContract, this.gameContract, response)
        this.currentState = this.agreedState
        return this.currentState
      })
    }

    return Promise.resolve(this.currentState)
  }

  async commit(move: Move): Promise<void> {
    if (this.pendingMoves[move.playerID] !== undefined) {
      throw Error(`player ${move.playerID} already committed`)
    }

    const state = await this.state

    if (move.playerID === this.playerID) {
      const response = await state.isMoveLegal(move)

      if (!response.isLegal) {
        throw Error(`illegal move: ${response.reason}`)
      }

      await move.sign(this.subkey, state)

    } else {
      const opponent = await this.arcadeumContract.playerAccount(this.timestamp, this.opponentTimestampSignature, this.opponentSubkeySignature)
      const moveMaker = await this.arcadeumContract.moveMaker(state.encoding, move, this.opponentSubkeySignature)

      if (moveMaker !== opponent) {
        throw Error(`move not signed by opponent`)
      }

      const response = await state.isMoveLegal(move)

      if (!response.isLegal) {
        if (await this.arcadeumContract.canReportCheater(this, state.encoding, move)) {
          await this.arcadeumContract.reportCheater(this, state.encoding, move)
        }

        throw Error(`illegal move: ${response.reason}`)
      }
    }

    const nextPlayers = await state.nextPlayers

    if (nextPlayers !== NextPlayers.Both) {
      if (move.playerID === this.playerID) {
        this.playerMoves.push(move)

      } else {
        this.agreedState = state
        this.opponentMove = move
        this.playerMoves = []
      }

      if (this.onCommit !== undefined) {
        this.onCommit(this, state, move)
      }

      this.currentState = await state.nextState(move)

      if (this.onChange !== undefined) {
        this.onChange(this, state, this.currentState, move)
      }

    } else {
      this.pendingMoves[move.playerID] = move

      if (this.pendingMoves[0] === undefined || this.pendingMoves[1] === undefined) {
        if (this.onCommit !== undefined) {
          this.onCommit(this, state, move)
        }

        return
      }

      this.agreedState = state
      this.opponentMove = this.pendingMoves[1 - this.playerID]
      this.playerMoves = [this.pendingMoves[this.playerID]!]

      if (this.onCommit !== undefined) {
        this.onCommit(this, state, move)
      }

      const pendingMoves = this.pendingMoves
      this.currentState = await state.nextState(this.pendingMoves[0]!, this.pendingMoves[1]!)
      this.pendingMoves = [undefined, undefined]

      if (this.onChange !== undefined) {
        this.onChange(this, state, this.currentState, pendingMoves[0]!, pendingMoves[1]!)
      }
    }

    const winner = await this.currentState.winner

    if (winner === Winner.Player0 && this.playerID === 0 || winner === Winner.Player1 && this.playerID === 1) {
      if (await this.arcadeumContract.canClaimReward(this, this.agreedState!.encoding, this.opponentMove, this.playerMoves)) {
        await this.arcadeumContract.claimReward(this, this.agreedState!.encoding, this.opponentMove, this.playerMoves)
      }
    }
  }

  readonly game: string
  readonly timestamp: ethers.utils.BigNumber
  readonly playerID: number
  readonly players: [PlayerInterface, PlayerInterface]
  readonly matchSignature: Signature
  readonly opponentSubkeySignature: Signature

  private get opponentID(): number {
    return 1 - this.playerID
  }

  private get opponentTimestampSignature(): Signature {
    return this.players[this.opponentID].timestampSignature
  }

  private agreedState?: State
  private opponentMove?: Move
  private playerMoves: Move[]
  private currentState?: State
  private pendingMoves: [Move | undefined, Move | undefined]
}

export class State {
  constructor(private arcadeumContract: ethers.Contract, protected gameContract: ethers.Contract, state: StateInterface) {
    this.tag = state.state.tag
    this.data = state.state.data
    this.metadata = {
      nonce: state.nonce,
      tag: state.tag,
      data: state.data
    }
  }

  get winner(): Promise<Winner> {
    return this.gameContract.winner(this.encoding)
  }

  get nextPlayers(): Promise<NextPlayers> {
    return this.gameContract.nextPlayers(this.encoding)
  }

  async isMoveLegal(move: Move): Promise<{ isLegal: boolean, reason: number }> {
    const response = await this.gameContract.isMoveLegal(this.encoding, move)

    return {
      isLegal: response[0],
      reason: response[1]
    }
  }

  async nextState(aMove: Move): Promise<State>
  async nextState(aMove: Move, anotherMove: Move): Promise<State>
  async nextState(aMove: [Move]): Promise<State>
  async nextState(aMove: [Move, Move]): Promise<State>
  async nextState(aMove: Move | [Move] | [Move, Move], anotherMove?: Move): Promise<State> {
    if (aMove instanceof Array) {
      if (anotherMove !== undefined) {
        throw Error(`unexpected second argument: array already given`)
      }

      switch (aMove.length) {
      case 1:
        return new State(this.arcadeumContract, this.gameContract, await this.gameContract[`nextState((uint32,uint8,bytes32[3],(uint32,bytes32[1])),(uint8,bytes))`](this.encoding, aMove[0]))

      case 2:
        return new State(this.arcadeumContract, this.gameContract, await this.gameContract[`nextState((uint32,uint8,bytes32[3],(uint32,bytes32[1])),(uint8,bytes),(uint8,bytes))`](this.encoding, aMove[0], aMove[1]))
      }

    } else {
      if (anotherMove === undefined) {
        return new State(this.arcadeumContract, this.gameContract, await this.gameContract[`nextState((uint32,uint8,bytes32[3],(uint32,bytes32[1])),(uint8,bytes))`](this.encoding, aMove))

      } else {
        return new State(this.arcadeumContract, this.gameContract, await this.gameContract[`nextState((uint32,uint8,bytes32[3],(uint32,bytes32[1])),(uint8,bytes),(uint8,bytes))`](this.encoding, aMove, anotherMove))
      }
    }

    throw Error(`expected dgame.Move[] of length 1 or 2`)
  }

  get encoding(): StateInterface {
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

  get hash(): Promise<ethers.utils.BigNumber> {
    return this.arcadeumContract.stateHash(this.encoding)
  }

  private tag: number
  // XXX: https://github.com/ethereum/solidity/issues/3270
  private data: [ethers.utils.BigNumber]
  private metadata: {
    nonce: number
    tag: number
    // XXX: https://github.com/ethereum/solidity/issues/3270
    data: [ethers.utils.BigNumber, ethers.utils.BigNumber, ethers.utils.BigNumber]
  }
}

export class Move {
  constructor(readonly move: { playerID: number, data: Uint8Array, signature?: Signature }) {
    this.playerID = move.playerID
    this.data = move.data

    if (move.signature !== undefined) {
      this.signature = move.signature
    } else {
      this.signature = new Signature()
    }
  }

  async sign(subkey: ethers.Wallet, state: State): Promise<void> {
    this.signature = sign(subkey, [`bytes32`, `uint8`, `bytes`], [await state.hash, this.playerID, this.data])
  }

  readonly playerID: number
  readonly data: Uint8Array
  signature: Signature
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
  readonly publicSeed: [ethers.utils.BigNumber]
  readonly timestampSignature: Signature
}

interface StateInterface {
  readonly nonce: number
  readonly tag: number
  // XXX: https://github.com/ethereum/solidity/issues/3270
  readonly data: [ethers.utils.BigNumber, ethers.utils.BigNumber, ethers.utils.BigNumber]
  readonly state: {
    readonly tag: number
    // XXX: https://github.com/ethereum/solidity/issues/3270
    readonly data: [ethers.utils.BigNumber]
  }
}

interface ChangeCallback {
  (match: Match, previousState: State, currentState: State, aMove: Move, anotherMove?: Move): void
}

interface CommitCallback {
  (match: Match, previousState: State, move: Move): void
}

class RemoteMatch extends Match {
  constructor(private relay: wsrelay.Relay, arcadeumContract: ethers.Contract, gameContract: ethers.Contract, subkey: ethers.Wallet, match: MatchInterface, onChange?: ChangeCallback, onCommit?: CommitCallback) {
    super(arcadeumContract, gameContract, subkey, match, onChange, onCommit)

    relay.subscribe(this)
  }

  async commit(move: Move): Promise<void> {
    await super.commit(move)

    if (move.playerID === this.playerID) {
      this.relay.send(JSON.stringify(move))
    }
  }

  async next(message: wsrelay.Message): Promise<void> {
    const move = JSON.parse(message.payload)

    move.data = deserializeUint8Array(move.data)
    move.signature.r = deserializeUint8Array(move.signature.r)
    move.signature.s = deserializeUint8Array(move.signature.s)

    return this.commit(new Move(move))
  }

  complete(): void {
  }

  error(error: any): void {
  }
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

function deserializeUint8Array(data: object): Uint8Array {
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
