# Arcadeum

## Network

A player enters the network by staking a security deposit into the Arcadeum contract that exists on-chain.
This deposit is locked until the player decides to stop playing games, which is when they can recover their original deposit.
The player is disincentivized from cheating because their deposit will be forfeit and no longer recoverable.

Players can request matches from a server that exists to facilitate the matchmaking process as well as relay messages between players.
The request begins with a declaration that they will use some randomly generated (ephemeral) key-pair for the duration of the next match ("I will use the key 0x1234... for the next match"), signed by the player's primary account that they used to stake their security deposit.
Using an ephemeral key-pair allows players to cryptographically sign moves without needing their primary account, while allowing others to confirm that those moves were indeed made by a given player.

In that same request, the player also declares what game they want to play, and a secret seed to initialize the state of the match with (this might be a deck of cards in a trading card game, for example).
This information currently isn't signed, but should be in the future to prevent the server from lying about player requests.
The server confirms that the player's balance in the network is sufficient and either pairs them up with another suitable player who's already awaiting a match, or puts them into a queue to wait for another player to be paired with.

Once the server decides to match two players together, it sends back an expiration timestamp for each to sign.
Each player uses their ephemeral key to sign the timestamp, which is a commitment by the player to the server that they promise not to exit the network until the timestamp expires, or until the match associated with the timestamp is completed and reported on-chain ("I will not withdraw my deposit until this match has completed and been reported or the timestamp has expired").

A player can exit the network by publicly announcing in a transaction that they intend on exiting the network.
This initiates a challenge period that allows another party (most likely the server) to cancel the exit if they possess that player's signed commitment not to exit.
If an exit is successfully challenged, the player's exit is disallowed and the challenger is reimbursed for transactional costs from the player's deposit.
If the exit challenge period elapses without any successful challenge against it, the player at that point simply issues a transaction that unlocks and returns their initial deposit.

The main benefit of this scheme is to reduce the number of on-chain transactions to at most one per match (zero in the event of ties which may not even need to be reported).

After receiving back the signed time commitments from both players, the match can then begin.
The server sends both players a data payload that includes computed public seeds for each player determined by the on-chain game contract.
That data payload also includes each player's timestamp signatures, subkey signatures, and a signature of the match details by the matcher.
From that payload, both players have enough information to derive the initial state of the match.

## State channel

Players exchange moves with their opponents through state channels, relayed via the matching server.
The match being played is a deterministic state matchine, where the inputs are cryptographically signed moves that are applied to the current state.
For the game logic, there are five functions that any game in Arcadeum must implement:

- `initialState()`: the starting state computed from the public seeds of both players
- `winnerImplementation()`: the decided winner for a given state
- `nextPlayersImplementation()`: who is allowed to apply moves on the current state
- `isMoveLegalImplementation()`: if a move applies to a given state or not
- `nextStateImplementation()`: the state resulting from applying a move(s) to a state

These functions are all implemented as pure functions that don't mutate the state of the Ethereum virtual machine, so they cost nothing to run on a local EVM node or bytecode interpreter.

Metastates in Arcadeum are higher-level states that wrap game-specific state data with Arcadeum-specific state data.
They allow us to implement useful primitives like commitment schemes and coordinated random number generation that can be used by game developers without having to reinvent the wheel.
For example, a developer can flag a given state as being determined by random entropy by writing `randomize(state)`, and then using that entropy in the `onRandomize()` callback.

Arcadeum moves are also similar to metastates.
They consist of the game-specific move data and a signature by the move maker of the game-specific move data and the state it was applied to.

Arcadeum does the work of validating proofs of victory and fraud, consulting with the game logic contract for game-specific behaviour and validation.

Fraud proofs are relatively simple.
They consist of a given state, and a given move applied to that state signed by the perpetrator of the fraud.
If the signature confirms that the move was made by the opponent, and the game contract confirms the move is invalid for the given state, then that's enough for the Arcadeum contract to validate that proof and slash the cheater's deposit.

Victory proofs are similar, but additionally require multiple moves by the winner to transition to a terminal win state.
They consist of a given state, a given move applied to that state signed by the opponent, and a series of moves that transitions to a terminal state where the winner has indeed won according to the game logic.
Because the winner sends the proof via signed transaction, the only other signature required is the signature by the opponent to establish that both players agree on a consensus state.
Beyond that, any additional moves made by the winner don't require any signing since they're already signed in the main transaction.
Arcadeum checks every single move in the proof against the game logic to validate it against the state it was applied to, transitioning to successive states with each move, and at last confirming that the final resulting state is terminal, with the prover being the victor.
Finally, the Arcadeum contract invalidates the timestamp for that match, allowing both players to exit again, and defers to the game contract to determine what payout the winner should receive.

Proving victory or fraud on-chain can be prohibitively expensive.
Victory proofs are generally more expensive to compute on-chain than fraud proofs since victory proofs involve verifying the validity of individual moves.
By our estimates, proving victory in our reference implementation of tic-tac-toe costs ~300k gas.
At the current median gas price, this translates to a cost of 307863 * 2.5 Gwei * 10e-9 ETH/Gwei * 433.72 USD = 0.33 USD per victory proof.
One potential approach to deal with this is to defer to the matching server as a privileged, but accountable authority for confirming a proof off-chain.
Players can submit their proofs back to the matcher, which then calls restricted functions on the Arcadeum contract that decide the result of the match without fully proving them.
The matcher would then be reimbursed for the cost of the transaction from the winner's deposit.
This would open a challenge period for players to challenge decisions made by the matcher.
If challenged, the onus is on the matcher to report the valid proof they received.
If the matcher decided the result of the match without possessing a proof, the challenger would be able to slash the matcher for a bounty.
If the matcher successfully answers the challenge, the cost of the full proof would be deducted from the challenger's deposit.
Therefore, challenging a result in favour of your opponent should only be done if one is certain that the opponent could not have won.
