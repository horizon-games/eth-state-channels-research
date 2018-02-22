import * as ethers from 'ethers'

export class Match {
  constructor(public game: string, public matchID: number) {
    this.players = [new Player(this), new Player(this)]
    this.signature = new Signature(0, `0x0000000000000000000000000000000000000000000000000000000000000000`, `0x0000000000000000000000000000000000000000000000000000000000000000`)
  }

  players: [Player, Player]
  signature: Signature

  async createAndSignSubkey(playerID: number): Promise<void> {
    this.players[playerID].createSubkey()

    return this.players[playerID].signSubkey()
  }

  async isSubkeySigned(playerID: number): Promise<boolean> {
    return this.contract.isSubkeySigned(this, playerID)
  }

  get initialState(): Promise<MetaState> {
    return this.contract.initialState(this).then(result => Object.assign(new MetaState(this), result[0]))
  }

  get contract(): ethers.Contract {
    if (this._contract === undefined || this._contract.address !== this.game) {
      const provider = new ethers.providers.JsonRpcProvider(`http://localhost:9545`) // XXX
      const metadata = require(`../../build/contracts/DGame.json`)

      this._contract = new ethers.Contract(this.game, metadata.abi, provider)
    }

    return this._contract
  }

  private _contract: ethers.Contract
}

export class Player {
  constructor(private owner: Match, public account?: string) {
    this.publicSeed = `0x`
  }

  subkey: string
  subkeySignature: Signature
  publicSeed: string

  createSubkey(): void {
    this.wallet = ethers.Wallet.createRandom()
    this.subkey = this.wallet.getAddress().toLowerCase()
  }

  async signSubkey(): Promise<void> {
    const provider = new ethers.providers.Web3Provider((global as any).web3.currentProvider)
    const accounts = await provider.listAccounts()

    const matchIDHex = ethers.utils.hexlify(this.owner.matchID)
    const matchIDBytes = ethers.utils.arrayify(matchIDHex)
    const matchIDPadded = ethers.utils.padZeros(matchIDBytes, 4)
    const matchIDPaddedHex = ethers.utils.hexlify(matchIDPadded)
    const message = `Sign to play! This won't cost anything.\n\nGame: ${this.owner.game}\nMatch: ${matchIDPaddedHex}\nPlayer: ${this.subkey}`

    this.account = accounts[0]
    this.subkeySignature = Signature.parse(await provider.getSigner(accounts[0]).signMessage(message))
  }

  private wallet: ethers.Wallet
}

export class MetaState {
  nonce: number
  tag: number
  data: string[]
  state: State

  constructor(private owner: Match) {
  }

  get winner(): Promise<number> {
    return this.contract.winner(this).then(result => result[0].toNumber())
  }

  get nextPlayers(): Promise<number> {
    return this.contract.nextPlayers(this).then(result => result[0].toNumber())
  }

  async isMoveLegal(move: Move): Promise<boolean> {
    return this.contract.isMoveLegal(this, move)
  }

  async nextState(a: Move | [Move] | [Move, Move], b?: Move): Promise<MetaState> {
    // XXX: https://github.com/ethers-io/ethers.js/issues/119
    if (a instanceof Move) {
      if (b === undefined) {
        return Object.assign(new MetaState(this.owner), (await this.contract.nextState1(this, a))[0])
      } else {
        return Object.assign(new MetaState(this.owner), (await this.contract.nextState2(this, a, b))[0])
      }
    }

    if (b !== undefined) {
      throw new Error(`invalid argument: ${b}`)
    }

    // XXX: https://github.com/ethers-io/ethers.js/issues/119
    switch (a.length) {
    case 1:
      return Object.assign(new MetaState(this.owner), (await this.contract.nextState1(this, a[0]))[0])
    case 2:
      return Object.assign(new MetaState(this.owner), (await this.contract.nextState2(this, a[0], a[1]))[0])
    }

    throw new Error(`expected 1 or 2 moves`)
  }

  get contract(): ethers.Contract {
    return this.owner.contract
  }
}

export class MetaMove {
  move: Move
  signature: Signature
}

export class Signature {
  constructor(public v: number, public r: string, public s: string) {
  }

  static parse(signature: string): Signature {
    const signatureBytes = ethers.utils.arrayify(signature)
    const r = ethers.utils.hexlify(new Uint8Array(signatureBytes.buffer, 0, 32))
    const s = ethers.utils.hexlify(new Uint8Array(signatureBytes.buffer, 32, 32))

    return new Signature(signatureBytes[64], r, s)
  }
}

export class State {
  tag: number
  data: string[]
}

export class Move {
  constructor(public playerID: number, public data: string) {
  }
}
