pragma solidity ^0.4.19;
pragma experimental ABIEncoderV2;

import './DGame.sol';

contract TTT is DGame {
  function winnerInternal(State state) internal pure returns (uint) {
    if (byte(0) != state.data[0] && state.data[0] == state.data[1] && state.data[1] == state.data[2]) {
      return uint(state.data[0]);

    } else if (byte(0) != state.data[3] && state.data[3] == state.data[4] && state.data[4] == state.data[5]) {
      return uint(state.data[3]);

    } else if (byte(0) != state.data[6] && state.data[6] == state.data[7] && state.data[7] == state.data[8]) {
      return uint(state.data[6]);

    } else if (byte(0) != state.data[0] && state.data[0] == state.data[3] && state.data[3] == state.data[6]) {
      return uint(state.data[0]);

    } else if (byte(0) != state.data[1] && state.data[1] == state.data[4] && state.data[4] == state.data[7]) {
      return uint(state.data[1]);

    } else if (byte(0) != state.data[2] && state.data[2] == state.data[5] && state.data[5] == state.data[8]) {
      return uint(state.data[2]);

    } else if (byte(0) != state.data[0] && state.data[0] == state.data[4] && state.data[4] == state.data[8]) {
      return uint(state.data[0]);

    } else if (byte(0) != state.data[2] && state.data[2] == state.data[4] && state.data[4] == state.data[6]) {
      return uint(state.data[2]);

    }

    return 0;
  }

  function nextPlayersInternal(State state) internal pure returns (uint) {
    if (state.tag >= 9) {
      return 0;
    }

    return state.tag % 2 + 1;
  }

  function isMoveLegalInternal(State state, Move move) internal pure returns (bool) {
    if (move.data[0] >= 9) {
      return false;
    }

    if (state.data[uint(move.data[0])] != 0) {
      return false;
    }

    return true;
  }

  function nextStateInternal(State state, Move move) internal pure returns (MetaState) {
    bytes memory data;

    data = new bytes(9);
    data[0] = state.data[0];
    data[1] = state.data[1];
    data[2] = state.data[2];
    data[3] = state.data[3];
    data[4] = state.data[4];
    data[5] = state.data[5];
    data[6] = state.data[6];
    data[7] = state.data[7];
    data[8] = state.data[8];
    data[uint(move.data[0])] = byte(move.playerID + 1);

    return play(State(state.tag + 1, data));
  }
}
