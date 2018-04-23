pragma solidity ^0.4.23;
pragma experimental ABIEncoderV2;

import 'arcadeum-contracts/contracts/DGame.sol';

contract TTT is DGame {
  uint32 private constant REASON_WRONG_LENGTH = 1;
  uint32 private constant REASON_NOT_A_SQUARE = 2;
  uint32 private constant REASON_ALREADY_PLAYED = 3;

  constructor(address owner) DGame(owner) public {
  }

  function matchDuration() public pure returns (uint) {
    return 10 minutes;
  }

  function initialState(bytes32[PUBLIC_SEED_LENGTH], bytes32[PUBLIC_SEED_LENGTH]) public pure returns (MetaState) {
    State memory state;

    return meta(state);
  }

  function winnerImplementation(State state) internal pure returns (Winner) {
    // XXX: https://github.com/ethereum/solidity/issues/3270
    if (byte(0) != state.data[0][0] && state.data[0][0] == state.data[0][1] && state.data[0][1] == state.data[0][2]) {
      return Winner(uint8(state.data[0][0]));

    // XXX: https://github.com/ethereum/solidity/issues/3270
    } else if (byte(0) != state.data[0][3] && state.data[0][3] == state.data[0][4] && state.data[0][4] == state.data[0][5]) {
      return Winner(uint8(state.data[0][3]));

    // XXX: https://github.com/ethereum/solidity/issues/3270
    } else if (byte(0) != state.data[0][6] && state.data[0][6] == state.data[0][7] && state.data[0][7] == state.data[0][8]) {
      return Winner(uint8(state.data[0][6]));

    // XXX: https://github.com/ethereum/solidity/issues/3270
    } else if (byte(0) != state.data[0][0] && state.data[0][0] == state.data[0][3] && state.data[0][3] == state.data[0][6]) {
      return Winner(uint8(state.data[0][0]));

    // XXX: https://github.com/ethereum/solidity/issues/3270
    } else if (byte(0) != state.data[0][1] && state.data[0][1] == state.data[0][4] && state.data[0][4] == state.data[0][7]) {
      return Winner(uint8(state.data[0][1]));

    // XXX: https://github.com/ethereum/solidity/issues/3270
    } else if (byte(0) != state.data[0][2] && state.data[0][2] == state.data[0][5] && state.data[0][5] == state.data[0][8]) {
      return Winner(uint8(state.data[0][2]));

    // XXX: https://github.com/ethereum/solidity/issues/3270
    } else if (byte(0) != state.data[0][0] && state.data[0][0] == state.data[0][4] && state.data[0][4] == state.data[0][8]) {
      return Winner(uint8(state.data[0][0]));

    // XXX: https://github.com/ethereum/solidity/issues/3270
    } else if (byte(0) != state.data[0][2] && state.data[0][2] == state.data[0][4] && state.data[0][4] == state.data[0][6]) {
      return Winner(uint8(state.data[0][2]));

    } else {
      return Winner.NONE;
    }
  }

  function nextPlayersImplementation(State state) internal pure returns (NextPlayers) {
    if (state.tag >= 9) {
      return NextPlayers.NONE;
    }

    return NextPlayers(1 + state.tag % 2);
  }

  function isMoveLegalImplementation(State state, Move move) internal pure returns (bool, uint32) {
    if (move.data.length != 1) {
      return (false, REASON_WRONG_LENGTH);
    }

    if (move.data[0] >= 9) {
      return (false, REASON_NOT_A_SQUARE);
    }

    // XXX: https://github.com/ethereum/solidity/issues/3270
    if (state.data[0][uint(move.data[0])] != 0) {
      return (false, REASON_ALREADY_PLAYED);
    }

    return (true, REASON_NONE);
  }

  function nextStateImplementation(State state, Move move) internal pure returns (MetaState) {
    State memory next;

    next.tag = state.tag + 1;
    // XXX: https://github.com/ethereum/solidity/issues/3270
    next.data[0] = state.data[0] | (bytes32(1 + move.playerID) << ((31 - uint(move.data[0])) * 8));

    return meta(next);
  }
}
