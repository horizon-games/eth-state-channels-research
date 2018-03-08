import * as ethers from 'ethers'

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
  constructor(gameAddress: string, private server: Server) {
    const arcadeumAddress = `0xcfeb869f69431e42cdb54a4f4f105c19c080a601`
    const arcadeumMetadata = require(`../../build/contracts/Arcadeum.json`)
    const gameMetadata = require(`../../build/contracts/DGame.json`)
    const provider = new ethers.providers.Web3Provider((window as any).web3.currentProvider)

    this.signer = provider.getSigner() // XXX: choose account
    this.arcadeumContract = new ethers.Contract(arcadeumAddress, arcadeumMetadata.abi, this.signer)
    this.gameContract = new ethers.Contract(gameAddress, gameMetadata.abi, this.signer)
  }

  get address(): string {
    return this.gameContract.address
  }

  get matchDuration(): Promise<number> {
    return this.gameContract.matchDuration().then(response => response.toNumber())
  }

  async isSecretSeedValid(address: string, secretSeed: Uint8Array): Promise<boolean> {
    return this.gameContract.isSecretSeedValid(address, secretSeed)
  }

  async createMatch(secretSeed: Uint8Array): Promise<{ match: Match, XXX: { subkeySignature: Signature, timestampSignature: Signature } }> {
    const subkey = ethers.Wallet.createRandom()
    const subkeyMessage = await this.arcadeumContract.subkeyMessage(subkey.getAddress())
    const subkeySignature = new Signature(await this.signer.signMessage(subkeyMessage))
    const timestamp = await this.server.sendSecretSeed(subkey.address, subkeySignature, secretSeed)
    const timestampSignature = sign(subkey, [`address`, `uint32`, `uint`], [this.address, timestamp.matchID, timestamp.timestamp])

    return {
      match: new Match(this.arcadeumContract, this.gameContract, subkey, await this.server.sendTimestampSignature(timestampSignature)),
      XXX: {
        subkeySignature: subkeySignature,
        timestampSignature: timestampSignature
      }
    }
  }

  private signer: ethers.providers.Web3Signer
  private arcadeumContract: ethers.Contract
  private gameContract: ethers.Contract
}

export interface Server {
  sendSecretSeed(subkeyAddress: string, subkeySignature: Signature, secretSeed: Uint8Array): Promise<TimestampInterface>
  sendTimestampSignature(timestampSignature: Signature): Promise<MatchInterface>
}

export class Match {
  constructor(private arcadeumContract: ethers.Contract, private gameContract: ethers.Contract, private subkey: ethers.Wallet, match: MatchInterface) {
    this.game = match.game
    this.matchID = match.matchID
    this.timestamp = match.timestamp
    this.playerID = match.playerID
    this.players = match.players
    this.matchSignature = match.matchSignature
    this.opponentSubkeySignature = match.opponentSubkeySignature
    this.opponentTimestampSignature = match.opponentTimestampSignature
    this.playerMoves = []
    this.pendingMoves = [undefined, undefined]
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
      const opponent = await this.arcadeumContract.playerAccount(this.game, this.matchID, this.timestamp, this.opponentTimestampSignature, this.opponentSubkeySignature)
      const moveMaker = await this.arcadeumContract.moveMaker(state.encoding, move, this.opponentSubkeySignature)

      if (moveMaker !== opponent) {
        throw Error(`move not signed by opponent`)
      }

      const response = await state.isMoveLegal(move)

      if (!response.isLegal) {
        this[`[object Object]`] = this.players // XXX

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

      this.currentState = await state.nextState(move)

    } else {
      this.pendingMoves[move.playerID] = move

      if (this.pendingMoves[0] === undefined || this.pendingMoves[1] === undefined) {
        return
      }

      this.agreedState = state
      this.opponentMove = this.pendingMoves[1 - this.playerID]
      this.playerMoves = [this.pendingMoves[this.playerID]!]
      this.currentState = await state.nextState(this.pendingMoves[0]!, this.pendingMoves[1]!)
      this.pendingMoves = [undefined, undefined]
    }

    const winner = await this.currentState.winner

    if (winner === Winner.Player0 && this.playerID === 0 || winner === Winner.Player1 && this.playerID === 1) {
      this[`[object Object]`] = this.players // XXX

      if (await this.arcadeumContract.canClaimReward(this, this.agreedState!.encoding, this.opponentMove, this.playerMoves)) {
        await this.arcadeumContract.claimReward(this, this.agreedState!.encoding, this.opponentMove, this.playerMoves)
      }
    }
  }

  readonly game: string
  readonly matchID: number
  readonly timestamp: ethers.utils.BigNumber
  readonly playerID: number
  readonly players: [PlayerInterface, PlayerInterface]
  readonly matchSignature: Signature
  readonly opponentSubkeySignature: Signature
  readonly opponentTimestampSignature: Signature

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
  constructor(readonly move: { playerID: number, data: Uint8Array }) {
    this.playerID = move.playerID
    this.data = move.data
    this.signature = new Signature()
  }

  async sign(subkey: ethers.Wallet, state: State): Promise<void> {
    this.signature = sign(subkey, [`bytes32`, `uint8`, `bytes`], [await state.hash, this.playerID, this.data])
  }

  readonly playerID: number
  readonly data: Uint8Array
  signature: Signature
}

export class Signature {
  constructor(signature?: string) {
    if (typeof signature === `string`) {
      const signatureBytes = ethers.utils.arrayify(signature)

      this.v = signatureBytes[64]
      this.r = new Uint8Array(signatureBytes.buffer, 0, 32)
      this.s = new Uint8Array(signatureBytes.buffer, 32, 32)

    } else {
      this.v = 0
      this.r = new Uint8Array(32)
      this.s = new Uint8Array(32)
    }
  }

  readonly v: number
  readonly r: Uint8Array
  readonly s: Uint8Array
}

interface TimestampInterface {
  readonly matchID: number
  readonly timestamp: ethers.utils.BigNumber
}

export interface MatchInterface {
  readonly game: string
  readonly matchID: number
  readonly timestamp: ethers.utils.BigNumber
  readonly playerID: number
  readonly players: [PlayerInterface, PlayerInterface]
  readonly matchSignature: Signature
  readonly opponentSubkeySignature: Signature
  readonly opponentTimestampSignature: Signature
}

interface PlayerInterface {
  readonly seedRating: number
  // XXX: https://github.com/ethereum/solidity/issues/3270
  readonly publicSeed: [ethers.utils.BigNumber]
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

function sign(wallet: ethers.Wallet, types: string[], values: any[]): Signature {
  const hash = ethers.utils.solidityKeccak256(types, values)
  const signatureValues = new ethers.SigningKey(wallet.privateKey).signDigest(hash)

  return {
    v: 27 + signatureValues.recoveryParam,
    r: ethers.utils.padZeros(ethers.utils.arrayify(signatureValues.r), 32),
    s: ethers.utils.padZeros(ethers.utils.arrayify(signatureValues.s), 32)
  }
}
