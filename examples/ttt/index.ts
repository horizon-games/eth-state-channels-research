import * as dgame from 'dgame'
import * as ethers from 'ethers'
import * as wsrelay from 'wsrelay'

class Matcher {
  constructor(private game: string) {
  }

  async sendSecretSeed(subkeyAddress: string, subkeySignature: dgame.Signature, secretSeed: Uint8Array) {
    const seed64 = base64(secretSeed)
    const r64 = base64(subkeySignature.r)
    const s64 = base64(subkeySignature.s)

    this.relay = new wsrelay.Relay(`localhost`, 8000, false, seed64, new wsrelay.Signature(subkeySignature.v, r64, s64), subkeyAddress, 1)

    const response = JSON.parse((await this.relay.connectForTimestamp()).payload)

    this.matchID = response.matchID
    this.timestamp = response.timestamp

    return response
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

    return JSON.parse((await this.relay.connectForMatchVerified()).payload)
  }

  private relay: wsrelay.Relay
  private matchID: number
  private timestamp: number
}

async function main(): Promise<void> {
  const gameAddress = `0xc89ce4735882c9f0f0fe26686c53074e09b0d550`
  const ttt = new dgame.DGame(gameAddress, new Matcher(gameAddress))

  await ttt.deposit(ethers.utils.parseEther(`1`))

  console.log(await ttt.matchDuration)
  console.log(await ttt.isSecretSeedValid(`0x0123456789012345678901234567890123456789`, new Uint8Array(0)))

  const match = (await ttt.createMatch(new Uint8Array(0))).match

  console.log(match)
}

function base64(data: Uint8Array): string {
  return new Buffer(ethers.utils.hexlify(data)).toString(`base64`)
}

main()
