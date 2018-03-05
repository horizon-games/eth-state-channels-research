import * as dgame from 'dgame'
import * as ethers from 'ethers'

main()

async function main(): Promise<void> {
  const ttt = new dgame.DGame(`0x8f0483125fcb9aaaefa9209d8e9d7b9c8b9fb90f`)

  console.log(await ttt.matchDuration)
  console.log(await ttt.isSecretSeedValid(`0x0123456789012345678901234567890123456789`, new Uint8Array(0)))

  const {
    match: match0,
    XXX: {
      subkeySignature: subkeySignature0,
      timestampSignature: timestampSignature0
    }
  } = await ttt.createMatch(new Uint8Array(0))

  const {
    match: match1,
    XXX: {
      subkeySignature: subkeySignature1,
      timestampSignature: timestampSignature1
    }
  } = await ttt.createMatch(new Uint8Array(0))

  match0.opponentSubkeySignature = subkeySignature1
  match0.opponentTimestampSignature = timestampSignature1

  match1.playerID = 1
  match1.opponentSubkeySignature = subkeySignature0
  match1.opponentTimestampSignature = timestampSignature0

  const provider = new ethers.providers.Web3Provider((window as any).web3.currentProvider)
  const accounts = await provider.listAccounts()
  const arcadeumAddress = `0x345ca3e014aaf5dca488057592ee47305d9b3e10`
  const arcadeumMetadata = require(`../../build/contracts/Arcadeum.json`)
  const arcadeumContract = new ethers.Contract(arcadeumAddress, arcadeumMetadata.abi, provider)
  const matchKey = new ethers.Wallet(`0xc87509a1c067bbde78beb793e6fa76530b6382a4c0241e5e4a9ec0a0f44dc0d3`, provider)
  const matchHash = await arcadeumContract.matchHash(match0.game, match0.matchID, match0.timestamp, [accounts[0], accounts[0]], [0, 0], [[`0x0000000000000000000000000000000000000000000000000000000000000000`], [`0x0000000000000000000000000000000000000000000000000000000000000000`]])
  const matchSignatureValues = new ethers.SigningKey(matchKey.privateKey).signDigest(matchHash)
  const matchSignature = {
    v: 27 + matchSignatureValues.recoveryParam,
    r: ethers.utils.padZeros(ethers.utils.arrayify(matchSignatureValues.r), 32),
    s: ethers.utils.padZeros(ethers.utils.arrayify(matchSignatureValues.s), 32)
  }

  match0.matchSignature = matchSignature
  match1.matchSignature = matchSignature

  console.log(match0)
  console.log(match1)

  const p = async (sender: dgame.Match, receiver: dgame.Match, square: number) => {
    const move = new dgame.Move({
      playerID: sender.playerID,
      data: new Uint8Array([square])
    })

    console.log(move)

    await sender.commit(move)
    await receiver.commit(move)

    console.log(sender)
    console.log(receiver)
  }

  const p0 = async (square: number) => p(match0, match1, square)
  const p1 = async (square: number) => p(match1, match0, square)

  await p0(0)
  await p1(4)
  await p0(8)
  await p1(2)
  await p0(6)
  await p1(3)
  await p0(7)
}
