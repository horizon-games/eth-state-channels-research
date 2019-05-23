# Ethereum State Channel Research

Enclosed is a collection of Solidity smart contracts as well as a TypeScript library for verifying turn-based games implemented as finite state machines in Solidity. This was in fact the original Arcadeum provably-fair game logic design, which includes a fully working tic-tac-toe game as a state channel. We've since progressed beyond this research to focus on WASM-based blockchain like EWASM, Substrate and others. Stay tuned for future publications of our latest research and we hope you enjoy this repo. However, this research is still relevent and can be utilized to build a token transfer design or other simple state updates. Through our experimentation we determined writing complex game logic in Solidity in a channel design isn't ideal and although it's possible for basic games, the constructions become complex with layers of machinery and its more important to have a game be cryptoeconomically secure than cryptographically secure. 


## Usage / Dev

**Tools:**

1. Install node v8.x or v10.x
2. `yarn install`
3. `yarn bootstrap`
4. `yarn build`


### Client

1. `cd client/`
2. `yarn test`
3. `yarn build`


### Ethereum

1. `cd ethereum/`
2. `yarn ganache` -- run in a separate terminal and import one of the private keys other than the first one into MetaMask
3. `yarn migrate`


### Server

1. Install redis-server and have it running in the background
2. `cd server/`
3. `make run`


### Running TTT example

1. `cd examples/ttt/`
2. `yarn install` -- this is required because `ttt` is configured outside of the arcadeum workspace
3. `yarn build` -- this will compile the TTT.sol ethereum contract
4. `yarn migrate` -- migrate the contract to the ganache testrpc
5. `yarn dev` -- this will start the webapp on http://localhost:3000/


### Building docs

Building the docs requires PlantUML and pdflatex.


# TODO, Build Optimizations

* ethereum/
  * dedicated abi.json files for each contract to reduce bundle filesize

* client/
  * optimize dist bundle filesize
  * potentially refactor out rxjs if turns out to be adding weight
