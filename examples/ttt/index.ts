import * as dgame from 'dgame'
import * as ethers from 'ethers'
import * as wsrelay from 'wsrelay'

class Matcher {
  constructor(private game: string) {
  }

  async sendSecretSeed(subkeyAddress: string, subkeySignature: dgame.Signature, secretSeed: Uint8Array): Promise<number> {
    const seed64 = base64(secretSeed)
    const r64 = base64(subkeySignature.r)
    const s64 = base64(subkeySignature.s)

    this.relay = new wsrelay.Relay(`localhost`, 8000, false, seed64, new wsrelay.Signature(subkeySignature.v, r64, s64), subkeyAddress, 1)
    this.relay.subscribe(this)

    const message = await this.relay.connectForTimestamp()

    this.matchID = message.meta.matchID
    this.timestamp = JSON.parse(message.payload)

    return this.timestamp
  }

  async sendTimestampSignature(timestampSignature: dgame.Signature): Promise<dgame.MatchInterface> {
    this.relay.send(JSON.stringify({
      gameID: 1,
      matchID: this.matchID,
      timestamp: this.timestamp,
      signature: {
        v: timestampSignature.v,
        r: base64(timestampSignature.r),
        s: base64(timestampSignature.s)
      }
    }), 2)

    const response = JSON.parse((await this.relay.connectForMatchVerified()).payload)

    response.players[0].publicSeed = [ethers.utils.bigNumberify(unbase64(response.players[0].publicSeed))]
    response.players[1].publicSeed = [ethers.utils.bigNumberify(unbase64(response.players[1].publicSeed))]
    response.matchSignature.r = unbase64(response.matchSignature.r)
    response.matchSignature.s = unbase64(response.matchSignature.s)
    response.opponentSubkeySignature.r = ethers.utils.arrayify(ethers.utils.toUtf8String(unbase64(response.opponentSubkeySignature.r)))
    response.opponentSubkeySignature.s = ethers.utils.arrayify(ethers.utils.toUtf8String(unbase64(response.opponentSubkeySignature.s)))
    response.opponentTimestampSignature.r = ethers.utils.arrayify(ethers.utils.toUtf8String(unbase64(response.opponentTimestampSignature.r)))
    response.opponentTimestampSignature.s = ethers.utils.arrayify(ethers.utils.toUtf8String(unbase64(response.opponentTimestampSignature.s)))

    return response
  }

  async commit(move: dgame.Move): Promise<void> {
    await this.match.commit(move)

    this.relay.send(JSON.stringify(move))
  }

  async next(message: wsrelay.Message): Promise<void> {
    const move = JSON.parse(message.payload)

    move.data = deserializeUint8Array(move.data)
    move.signature.r = deserializeUint8Array(move.signature.r)
    move.signature.s = deserializeUint8Array(move.signature.s)

    const metaMove = new dgame.Move(move)

    metaMove.signature = new dgame.Signature(move.signature)

    console.log(move)
    console.log(metaMove)

    await this.match.commit(metaMove)

    console.log(await this.match.state)

    switch (this.match.playerID) {
    case 0:
      switch (((await this.match.state) as any).tag) {
      case 2:
        await this.commit(new dgame.Move({ playerID: 0, data: new Uint8Array([8]) }))
        break

      case 4:
        await this.commit(new dgame.Move({ playerID: 0, data: new Uint8Array([6]) }))
        break

      case 6:
        await this.commit(new dgame.Move({ playerID: 0, data: new Uint8Array([7]) }))
        break
      }

      break

    case 1:
      switch (((await this.match.state) as any).tag) {
      case 1:
        await this.commit(new dgame.Move({ playerID: 1, data: new Uint8Array([4]) }))
        break

      case 3:
        await this.commit(new dgame.Move({ playerID: 1, data: new Uint8Array([2]) }))
        break

      case 5:
        await this.commit(new dgame.Move({ playerID: 1, data: new Uint8Array([3]) }))
        break
      }

      break
    }

    console.log(await this.match.state)
  }

  error(error: any): void {
  }

  complete(): void {
  }

  match: dgame.Match

  private relay: wsrelay.Relay
  private matchID: number
  private timestamp: number
}

async function main(): Promise<void> {
  const gameAddress = `0xc89ce4735882c9f0f0fe26686c53074e09b0d550`
  const matcher = new Matcher(gameAddress)
  const ttt = new dgame.DGame(gameAddress, matcher)

  await ttt.deposit(ethers.utils.parseEther(`1`))

  console.log(await ttt.matchDuration)
  console.log(await ttt.isSecretSeedValid(`0x0123456789012345678901234567890123456789`, new Uint8Array(0)))

  matcher.match = await ttt.createMatch(new Uint8Array(0))

  console.log(matcher.match)

  if (matcher.match.playerID === 0) {
    await matcher.commit(new dgame.Move({ playerID: 0, data: new Uint8Array([0]) }))
  }
}

function base64(data: Uint8Array): string {
  return new Buffer(ethers.utils.hexlify(data)).toString(`base64`)
}

function unbase64(data: string): Uint8Array {
  return Uint8Array.from(Buffer.from(data, `base64`))
}

function deserializeUint8Array(data: object): Uint8Array {
  const array: number[] = []

  for (let i = 0; data[i] !== undefined; i++) {
    array.push(data[i])
  }

  return new Uint8Array(array)
}

main()
