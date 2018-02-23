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
  console.log(state)
  console.log(`winner: ${await state.winner}`)
  console.log(`next player(s): ${await state.nextPlayers}`)

  state = await state.nextState(new dgame.Move(0, `0x00`))
  console.log(state)
  console.log(`winner: ${await state.winner}`)
  console.log(`next player(s): ${await state.nextPlayers}`)

  state = await state.nextState(new dgame.Move(1, `0x04`))
  console.log(state)
  console.log(`winner: ${await state.winner}`)
  console.log(`next player(s): ${await state.nextPlayers}`)

  state = await state.nextState(new dgame.Move(0, `0x08`))
  console.log(state)
  console.log(`winner: ${await state.winner}`)
  console.log(`next player(s): ${await state.nextPlayers}`)

  state = await state.nextState(new dgame.Move(1, `0x02`))
  console.log(state)
  console.log(`winner: ${await state.winner}`)
  console.log(`next player(s): ${await state.nextPlayers}`)

  state = await state.nextState(new dgame.Move(0, `0x06`))
  console.log(state)
  console.log(`winner: ${await state.winner}`)
  console.log(`next player(s): ${await state.nextPlayers}`)

  state = await state.nextState(new dgame.Move(1, `0x03`))
  console.log(state)
  console.log(`winner: ${await state.winner}`)
  console.log(`next player(s): ${await state.nextPlayers}`)

  state = await state.nextState(new dgame.Move(0, `0x07`))
  console.log(state)
  console.log(`winner: ${await state.winner}`)
  console.log(`next player(s): ${await state.nextPlayers}`)
}
