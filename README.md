# Arcadeum

Welcome to the ARCADEUM.network


## Usage / Dev

**Tools:**

1. Install node v8.x or v9.x
2. `yarn install`
3. `yarn bootstrap`


### Client

1. cd client/
2. yarn test
3. yarn build


### Ethereum

1. cd ethereum/
2. yarn testrpc - in a separate terminal
3. yarn migrate


### Server

1. cd server/
2. 


# TODO, Build Optimizations

* ethereum/
  * dedicate abi.json files for each contract

* client/
  * why is filesize so large? minified at 544 kb for such a small code-base
  * analyze bundle size
  * potentially refactor out rxjs if turns out to be adding weight



# OLD:

Arcadeum is a collection of Solidity smart contracts as well as a TypeScript library for verifying turn-based games implemented as finite state machines in Solidity.

## Usage

```yarn```

```yarn build```

```yarn testrpc```

```yarn migrate```

```yarn server```

```yarn dev```

Open a browser with MetaMask installed and go to http://localhost:3000/.
