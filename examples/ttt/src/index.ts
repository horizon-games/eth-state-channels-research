import * as dgame from 'arcadeum'
import * as ethers from 'ethers'

const ttt = new dgame.DGame(`0xd833215cbcc3f914bd1c9ece3ee7bf8b14f841bb`)
const signer = (ttt as any).signer
const arcadeumContract = (ttt as any).arcadeumContract

async function deposit(): Promise<void> {
  return ttt.deposit(ethers.utils.parseEther(`1`))
}

async function startWithdrawal(): Promise<void> {
  return arcadeumContract.startWithdrawal()
}

async function finishWithdrawal(): Promise<void> {
  return arcadeumContract.finishWithdrawal()
}

async function createMatch(): Promise<void> {
  console.log(await ttt.matchDuration)
  console.log(await ttt.isSecretSeedValid(`0x0123456789012345678901234567890123456789`, new Uint8Array(0)))

  const match = await ttt.createMatch(new Uint8Array(0), (match: dgame.Match, previousState: dgame.State, currentState: dgame.State, aMove: dgame.Move, anotherMove?: dgame.Move) => {
    console.log(currentState)

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
  console.log(await match.state)

  if (match.playerID === 0) {
    return match.commit(new dgame.Move({ playerID: 0, data: new Uint8Array([0]) }))
  }
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
