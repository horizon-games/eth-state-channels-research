pragma solidity ^0.4.19;

import './DGame.sol';

contract TTT is DGame {
  function playerStatusInternal(State state, uint8 player) internal pure returns (Status) {
    uint8 winner;

    winner = result(state);

    if (winner != 0) {
      if (player + 1 == winner) {
        return Status.Won;
      } else {
        return Status.Done;
      }
    }

    if (state.type_ % 2 == player) {
      return Status.Moving;
    } else {
      return Status.Waiting;
    }
  }

  function isMoveLegalInternal(State state, Move move) internal pure returns (bool) {
    if (move.data.length != 1) {
      return false;
    }

    if (move.data[0] >= 9) {
      return false;
    }

    if (state.data[uint(move.data[0])] != 0) {
      return false;
    }

    return true;
  }

  function nextStateInternal(State state, Move[] moves) internal pure returns (DState) {
    bytes memory data;

    data = new bytes(state.data.length);
    data[0] = state.data[0];
    data[1] = state.data[1];
    data[2] = state.data[2];
    data[3] = state.data[3];
    data[4] = state.data[4];
    data[5] = state.data[5];
    data[6] = state.data[6];
    data[7] = state.data[7];
    data[8] = state.data[8];
    data[uint(moves[0].data[0])] = byte(moves[0].player);

    return id(State(state.type_ + 1, data, new Status[](0)));
  }

  function result(State state) private pure returns (uint8) {
    if (state.data.length == 0) {
      return 0;
    }

    if (byte(0) != state.data[0] && state.data[0] == state.data[1] && state.data[1] == state.data[2]) {
      return uint8(state.data[0]);
    } else if (byte(0) != state.data[3] && state.data[3] == state.data[4] && state.data[4] == state.data[5]) {
      return uint8(state.data[3]);
    } else if (byte(0) != state.data[6] && state.data[6] == state.data[7] && state.data[7] == state.data[8]) {
      return uint8(state.data[6]);
    } else if (byte(0) != state.data[0] && state.data[0] == state.data[3] && state.data[3] == state.data[6]) {
      return uint8(state.data[0]);
    } else if (byte(0) != state.data[1] && state.data[1] == state.data[4] && state.data[4] == state.data[7]) {
      return uint8(state.data[1]);
    } else if (byte(0) != state.data[2] && state.data[2] == state.data[5] && state.data[5] == state.data[8]) {
      return uint8(state.data[2]);
    } else if (byte(0) != state.data[0] && state.data[0] == state.data[4] && state.data[4] == state.data[8]) {
      return uint8(state.data[0]);
    } else if (byte(0) != state.data[2] && state.data[2] == state.data[4] && state.data[4] == state.data[6]) {
      return uint8(state.data[2]);
    }

    return 0;
  }
}
