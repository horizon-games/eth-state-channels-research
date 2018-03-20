import * as dgame from 'dgame'
import * as ethers from 'ethers'

async function main(): Promise<void> {
  const ttt = new dgame.DGame(`0xc89ce4735882c9f0f0fe26686c53074e09b0d550`)

  await ttt.deposit(ethers.utils.parseEther(`1`))

  console.log(await ttt.matchDuration)
  console.log(await ttt.isSecretSeedValid(`0x0123456789012345678901234567890123456789`, new Uint8Array(0)))

  const match = await ttt.createMatch(new Uint8Array(0), (match: dgame.Match, previousState: dgame.State, currentState: dgame.State, aMove: dgame.Move, anotherMove?: dgame.Move) => {
    switch (match.playerID) {
    case 0:
      switch ((currentState as any).tag) {
      case 2:
        match.commit(new dgame.Move({ playerID: 0, data: new Uint8Array([8]) }))
        break

      case 4:
        match.commit(new dgame.Move({ playerID: 0, data: new Uint8Array([6]) }))
        break

      case 6:
        match.commit(new dgame.Move({ playerID: 0, data: new Uint8Array([7]) }))
        break
      }

      break

    case 1:
      switch ((currentState as any).tag) {
      case 1:
        match.commit(new dgame.Move({ playerID: 1, data: new Uint8Array([4]) }))
        break

      case 3:
        match.commit(new dgame.Move({ playerID: 1, data: new Uint8Array([2]) }))
        break

      case 5:
        match.commit(new dgame.Move({ playerID: 1, data: new Uint8Array([3]) }))
        break
      }

      break
    }
  })

  console.log(match)

  if (match.playerID === 0) {
    return match.commit(new dgame.Move({ playerID: 0, data: new Uint8Array([0]) }))
  }
}

main()
