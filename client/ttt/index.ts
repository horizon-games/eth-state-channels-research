import * as dgame from 'dgame'

main()

async function main(): Promise<void> {
  const match = new dgame.Match(`0x8f0483125fcb9aaaefa9209d8e9d7b9c8b9fb90f`, 12345)
  match[`[object Object]`] = match.players // XXX
  await match.createAndSignSubkey(0)
  await match.createAndSignSubkey(1)

  if (await match.isSubkeySigned(0)) {
    console.log(`subkey 0 is signed`)
  } else {
    console.log(`subkey 0 not signed`)
  }

  if (await match.isSubkeySigned(1)) {
    console.log(`subkey 1 is signed`)
  } else {
    console.log(`subkey 1 not signed`)
  }

  let state = await match.initialState
  let winner = await state.winner
  let players = await state.nextPlayers
  console.log(state)
  console.log(`winner: ${winner}`)
  console.log(`next player(s): ${players}`)

  state = await state.nextState(new dgame.Move(0, `0x00`))
  winner = await state.winner
  players = await state.nextPlayers
  console.log(state)
  console.log(`winner: ${winner}`)
  console.log(`next player(s): ${players}`)

  state = await state.nextState(new dgame.Move(1, `0x04`))
  winner = await state.winner
  players = await state.nextPlayers
  console.log(state)
  console.log(`winner: ${winner}`)
  console.log(`next player(s): ${players}`)

  state = await state.nextState(new dgame.Move(0, `0x08`))
  winner = await state.winner
  players = await state.nextPlayers
  console.log(state)
  console.log(`winner: ${winner}`)
  console.log(`next player(s): ${players}`)

  state = await state.nextState(new dgame.Move(1, `0x02`))
  winner = await state.winner
  players = await state.nextPlayers
  console.log(state)
  console.log(`winner: ${winner}`)
  console.log(`next player(s): ${players}`)

  state = await state.nextState(new dgame.Move(0, `0x06`))
  winner = await state.winner
  players = await state.nextPlayers
  console.log(state)
  console.log(`winner: ${winner}`)
  console.log(`next player(s): ${players}`)

  state = await state.nextState(new dgame.Move(1, `0x03`))
  winner = await state.winner
  players = await state.nextPlayers
  console.log(state)
  console.log(`winner: ${winner}`)
  console.log(`next player(s): ${players}`)

  state = await state.nextState(new dgame.Move(0, `0x07`))
  winner = await state.winner
  players = await state.nextPlayers
  console.log(state)
  console.log(`winner: ${winner}`)
  console.log(`next player(s): ${players}`)
}
