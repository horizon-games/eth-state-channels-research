pragma solidity ^0.4.19;
pragma experimental ABIEncoderV2;

import './DGame.sol';

contract TTT is DGame {
  function TTT(Arcadeum arcadeum) DGame(arcadeum) public {
  }

  function winnerInternal(State state) internal pure returns (uint) {
    if (byte(0) != state.data[0][0] && state.data[0][0] == state.data[0][1] && state.data[0][1] == state.data[0][2]) {
      return uint(state.data[0][0]);

    } else if (byte(0) != state.data[0][3] && state.data[0][3] == state.data[0][4] && state.data[0][4] == state.data[0][5]) {
      return uint(state.data[0][3]);

    } else if (byte(0) != state.data[0][6] && state.data[0][6] == state.data[0][7] && state.data[0][7] == state.data[0][8]) {
      return uint(state.data[0][6]);

    } else if (byte(0) != state.data[0][0] && state.data[0][0] == state.data[0][3] && state.data[0][3] == state.data[0][6]) {
      return uint(state.data[0][0]);

    } else if (byte(0) != state.data[0][1] && state.data[0][1] == state.data[0][4] && state.data[0][4] == state.data[0][7]) {
      return uint(state.data[0][1]);

    } else if (byte(0) != state.data[0][2] && state.data[0][2] == state.data[0][5] && state.data[0][5] == state.data[0][8]) {
      return uint(state.data[0][2]);

    } else if (byte(0) != state.data[0][0] && state.data[0][0] == state.data[0][4] && state.data[0][4] == state.data[0][8]) {
      return uint(state.data[0][0]);

    } else if (byte(0) != state.data[0][2] && state.data[0][2] == state.data[0][4] && state.data[0][4] == state.data[0][6]) {
      return uint(state.data[0][2]);

    }

    return 0;
  }

  function nextPlayersInternal(State state) internal pure returns (uint) {
    if (state.tag >= 9) {
      return 0;
    }

    return 1 + state.tag % 2;
  }

  function isMoveLegalInternal(State state, Move move) internal pure returns (bool) {
    if (move.data[0] >= 9) {
      return false;
    }

    if (state.data[0][uint(move.data[0])] != 0) {
      return false;
    }

    return true;
  }

  function nextStateInternal(State state, Move move) internal pure returns (MetaState) {
    bytes32 mask;

    mask = bytes32(1 + move.playerID);
    mask <<= (31 - uint(move.data[0])) * 8;

    return play(State(state.tag + 1, [state.data[0] | mask]));
  }
}
