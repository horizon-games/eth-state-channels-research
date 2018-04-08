import * as arcadeum from 'arcadeum'
import * as ethers from 'ethers'

const game = new arcadeum.Game(`0xd833215cbcc3f914bd1c9ece3ee7bf8b14f841bb`)
const signer = (game as any).signer
const arcadeumContract = (game as any).arcadeumContract

async function deposit(): Promise<string> {
  return game.deposit(ethers.utils.bigNumberify(`1000000000000000000`))
}

async function startWithdrawal(): Promise<void> {
  return arcadeumContract.startWithdrawal()
}

async function finishWithdrawal(): Promise<void> {
  return arcadeumContract.finishWithdrawal()
}

async function createMatch(): Promise<void> {
  const match = game.createMatch(new Uint8Array(0))

  match.addCallback(async (nextState: arcadeum.State, previousState?: arcadeum.State, aMove?: arcadeum.Move, anotherMove?: arcadeum.Move) => {
    console.log(aMove)
    console.log(nextState)

    switch (match.playerID) {
    case 0:
      switch ((nextState as any).state.tag) {
      case 0:
        match.queueMove(await match.createMove(new Uint8Array([0])))
        break

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
      switch ((nextState as any).state.tag) {
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

  console.log(match)

  return match.ready
}

(window as any).deposit = deposit;
(window as any).startWithdrawal = startWithdrawal;
(window as any).finishWithdrawal = finishWithdrawal;
(window as any).createMatch = createMatch

window.setInterval(async () => {
  (document.getElementById(`currentTime`))!.textContent = new Date().toLocaleString()

  const account = await signer.getAddress()

  arcadeumContract.balance(account).then((balance) => {
    (document.getElementById(`balance`))!.textContent = balance.toString()
  })

  arcadeumContract.withdrawalTime(account).then((withdrawalTime) => {
    (document.getElementById(`withdrawalTime`))!.textContent = new Date(1000 * withdrawalTime.toNumber()).toLocaleString()
  })
}, 1000)
