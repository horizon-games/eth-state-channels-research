import * as arcadeum from 'arcadeum'
import * as ethers from 'ethers'

const game = new arcadeum.Game(`0x9561c133dd8580860b6b7e504bc5aa500f0f06a7`)
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

    render((nextState as any).metaState.state.data)
  })

  console.log(match)

  await match.ready

  window[`match`] = match
}

function render(squares: number[]): void {
  const canvas = document.getElementById(`canvas`) as HTMLCanvasElement
  const { width, height } = canvas
  const context = canvas.getContext(`2d`)

  if (context === null) {
    throw Error(`no canvas context`)
  }

  context.fillStyle = `white`
  context.fillRect(0, 0, width, height)

  context.strokeStyle = `black`
  context.lineWidth = 5
  context.lineCap = `round`
  context.lineJoin = `round`
  context.beginPath()
  context.moveTo((1 / 3) * width, 10)
  context.lineTo((1 / 3) * width, height - 10)
  context.moveTo((2 / 3) * width, 10)
  context.lineTo((2 / 3) * width, height - 10)
  context.moveTo(10, (1 / 3) * height)
  context.lineTo(width - 10, (1 / 3) * height)
  context.moveTo(10, (2 / 3) * height)
  context.lineTo(width - 10, (2 / 3) * height)
  context.stroke()

  for (let i = 0; i < 3; i++) {
    for (let j = 0; j < 3; j++) {
      switch (squares[3 * i + j]) {
      case 1:
        context.strokeStyle = `red`
        context.lineWidth = 10
        context.beginPath()
        context.moveTo((j / 3) * width + 30, (i / 3) * height + 30)
        context.lineTo(((j + 1) / 3) * width - 30, ((i + 1) / 3) * height - 30)
        context.moveTo(((j + 1) / 3) * width - 30, (i / 3) * height + 30)
        context.lineTo((j / 3) * width + 30, ((i + 1) / 3) * height - 30)
        context.stroke()
        break

      case 2:
        context.strokeStyle = `blue`
        context.lineWidth = 10
        context.beginPath()
        context.ellipse(((j + 0.5) / 3) * width, ((i + 0.5) / 3) * height, width / 6 - 30, height / 6 - 30, 0, 0, 2 * Math.PI)
        context.stroke()
        break
      }
    }
  }
}

window[`deposit`] = deposit
window[`startWithdrawal`] = startWithdrawal
window[`finishWithdrawal`] = finishWithdrawal
window[`createMatch`] = createMatch

let renderCanvas = true

setInterval(async () => {
  (document.getElementById(`currentTime`))!.textContent = new Date().toLocaleString()

  const account = await signer.getAddress()

  arcadeumContract.balance(account).then((balance) => {
    (document.getElementById(`balance`))!.textContent = balance.toString()
  })

  arcadeumContract.withdrawalTime(account).then((withdrawalTime) => {
    (document.getElementById(`withdrawalTime`))!.textContent = new Date(1000 * withdrawalTime.toNumber()).toLocaleString()
  })

  if (renderCanvas && document.getElementById(`canvas`)) {
    renderCanvas = false

    render([0, 0, 0, 0, 0, 0, 0, 0, 0])

    const canvas = document.getElementById(`canvas`) as HTMLCanvasElement
    canvas.addEventListener(`click`, async (event: MouseEvent) => {
      const match = window[`match`]
      const { width, height } = canvas

      if (match === undefined) {
        return
      }

      match.queueMove(await match.createMove(new Uint8Array([3 * Math.floor(event.offsetY / (height / 3)) + Math.floor(event.offsetX / (width / 3))])))
    })
  }
}, 1000)
