import * as dgame from 'arcadeum'
import * as ethers from 'ethers'
import { arcadeumAddress, gameAddress, serverAddress, deposit } from './arcadeum'
import { wallet1, wallet2 } from './wallet'

describe('ttt', () => {
  it('should successfully complete an end-to-end game', async (done) => {
    const ttt = new dgame.Game(gameAddress, { arcadeumAddress: arcadeumAddress, serverAddress: serverAddress, wallet: wallet1 })
    const ttt2 = new dgame.Game(gameAddress, { arcadeumAddress: arcadeumAddress, serverAddress: serverAddress, wallet: wallet2 })
    const arcadeumContract = (ttt as any).arcadeumContract
    const depositInWei = ethers.utils.parseEther(deposit)
    const balanceInWei = await arcadeumContract.balance(wallet1.address) as ethers.utils.BigNumber
    const balance2InWei = await arcadeumContract.balance(wallet2.address) as ethers.utils.BigNumber
    if (balanceInWei.lt(depositInWei)) {
      console.log(`staking ${deposit} ETH for wallet ${wallet1.address}`)
      const response = await ttt.deposit(depositInWei)

      if (wallet1.provider !== undefined) {
        const transaction = await wallet1.provider.waitForTransaction(response, 60000)
        console.log(`transaction hash mined ${transaction.hash}`)
      }
    }
    if (balance2InWei.lt(depositInWei)) {
      console.log(`staking ${deposit} ETH for wallet ${wallet2.address}`)
      const response = await ttt2.deposit(depositInWei)

      if (wallet2.provider !== undefined) {
        const transaction = await wallet2.provider.waitForTransaction(response, 60000)
        console.log(`transaction hash mined ${transaction.hash}`)
      }
    }
    console.log('begin match')
    Promise.all([createMatch(ttt), createMatch(ttt2)]).then(values => {
      console.log('Winner!')
      console.log(values)
      done()
    }).catch(e => {
      console.log('Error!')
      console.log(e)
      done(e)
    })
  }, 200000)
})

// Client game logic that would normally run in the browser
async function createMatch(game: dgame.Game): Promise<dgame.Winner> {
  return new Promise<dgame.Winner>(async (resolve, reject) => {
    const match = await game.createMatch(new Uint8Array(0))

    match.addCallback(async (nextState: dgame.State, previousState?: dgame.State, aMove?: dgame.Move, anotherMove?: dgame.Move) => {
      const winner = await nextState.winner
      if (winner !== dgame.Winner.None) {
        resolve(winner)
      }

      switch (match.playerID) {
        case 0:
          switch ((nextState as any).tag) {
            case 2:
              match.queueMove(await match.createMove(new Uint8Array([8])))
              break

            case 4:
              match.queueMove(await match.createMove(new Uint8Array([6])))
              break

            case 6:
              match.queueMove(await match.createMove(new Uint8Array([7])))
              break
          }

          break

        case 1:
          switch ((nextState as any).tag) {
            case 1:
              match.queueMove(await match.createMove(new Uint8Array([4])))
              break

            case 3:
              match.queueMove(await match.createMove(new Uint8Array([2])))
              break

            case 5:
              match.queueMove(await match.createMove(new Uint8Array([3])))
              break
          }

          break
      }
    })

    if (match.playerID === 0) {
      match.queueMove(await match.createMove(new Uint8Array([0])))
    }
  })
}
