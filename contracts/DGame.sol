pragma solidity ^0.4.19;

contract DGame {
  enum Status {
    Playing,
    Waiting,
    Moving,
    Done,
    Won,
    Cheated
  }

  enum DType {
    None,
    CommittingRandom,
    RevealingRandom,
    CommittingSecret,
    RevealingSecret
  }

  struct DState {
    DType type_;
    bytes data;
    Status[] statuses;
    State state;
  }

  struct State {
    uint8 type_;
    bytes data;
    Status[] statuses;
  }

  struct Move {
    uint8 player;
    bytes data;
  }

  struct Signature {
    uint8 v;
    bytes32 r;
    bytes32 s;
  }

  function playerStatus(State state, uint8 player) internal pure returns (Status);
  function isMoveLegal(State state, Move move) internal pure returns (bool);
  function nextState(State state, Move[] moves) internal pure returns (DState);

  function onRandomize(State state, bytes) internal pure returns (DState) {
    return id(state);
  }

  function onExchange(State state, Move[]) internal pure returns (DState) {
    return id(state);
  }

  function id(State state) internal pure returns (DState) {
    return DState(DType.None, new bytes(0), new Status[](0), state);
  }

  function randomize(State state, uint8 bytes_, uint8[] players) internal pure returns (DState) {
    bytes memory data;
    Status[] memory statuses;
    Status status;
    uint8 i;

    data = new bytes(1);
    data[0] = byte(bytes_);

    statuses = new Status[](state.statuses.length);

    for (i = 0; i < statuses.length; i++) {
      status = state.statuses[i];

      if (status == Status.Playing || status == Status.Moving) {
        status = Status.Waiting;
      }

      statuses[i] = status;
    }

    for (i = 0; i < players.length; i++) {
      status = statuses[players[i]];

      if (status == Status.Waiting) {
        statuses[players[i]] = Status.Moving;
      }
    }

    return DState(DType.CommittingRandom, data, statuses, state);
  }

  function exchange(State state, uint8[] players) internal pure returns (DState) {
    Status[] memory statuses;
    Status status;
    uint8 i;

    statuses = new Status[](state.statuses.length);

    for (i = 0; i < statuses.length; i++) {
      status = state.statuses[i];

      if (status == Status.Playing || status == Status.Moving) {
        status = Status.Waiting;
      }

      statuses[i] = status;
    }

    for (i = 0; i < players.length; i++) {
      status = statuses[players[i]];

      if (status == Status.Waiting) {
        statuses[players[i]] = Status.Moving;
      }
    }

    return DState(DType.CommittingSecret, new bytes(0), statuses, state);
  }

  function playerStatusInternal(DState state, uint8 player) private pure returns (Status) {
    Status status;

    status = Status.Playing;

    if (state.statuses.length > player) {
      status = state.statuses[player];
    }

    if (status == Status.Playing) {
      if (state.state.statuses.length > player) {
        status = state.state.statuses[player];
      }

      if (status == Status.Playing) {
        status = playerStatus(state.state, player);
      }
    }

    return status;
  }

  function isMoveLegalInternal(DState state, Move move) private pure returns (bool) {
    bytes32 hash;
    uint8 i;
    uint8 j;

    if (playerStatusInternal(state, move.player) != Status.Moving) {
      return false;
    }

    if (state.type_ == DType.None) {
      return isMoveLegal(state.state, move);
    } else if (state.type_ == DType.CommittingRandom) {
      return move.data.length == 32;
    } else if (state.type_ == DType.RevealingRandom) {
      if (move.data.length != uint(state.data[0])) {
        return false;
      }

      hash = keccak256(move.data);
      i = uint8(state.data[1 + move.player]);

      for (j = 0; j < 32; j++) {
        if (hash[j] != state.data[1 + state.statuses.length + i * 32 + j]) {
          return false;
        }
      }

      return true;
    } else if (state.type_ == DType.CommittingSecret) {
      return move.data.length == 32;
    } else if (state.type_ == DType.RevealingSecret) {
      hash = keccak256(move.data);
      i = uint8(state.data[move.player]);

      for (j = 0; j < 32; j++) {
        if (hash[j] != state.data[state.statuses.length + i * 32 + j]) {
          return false;
        }
      }

      return true;
    }
  }

  function nextStateInternal(DState state, Move[] moves) private pure returns (DState) {
    bytes memory data;
    uint8 i;
    uint8 j;

    if (state.type_ == DType.None) {
      return nextState(state.state, moves);
    } else if (state.type_ == DType.CommittingRandom) {
      data = new bytes(1 + state.statuses.length + moves.length * 32);
      data[0] = state.data[0];

      for (i = 0; i < moves.length; i++) {
        data[1 + moves[i].player] = byte(i);

        for (j = 0; j < 32; j++) {
          data[1 + state.statuses.length + i * 32 + j] = moves[i].data[j];
        }
      }

      return DState(DType.RevealingRandom, data, state.statuses, state.state);
    } else if (state.type_ == DType.RevealingRandom) {
      data = new bytes(uint(state.data[0]));

      for (i = 0; i < moves.length; i++) {
        for (j = 0; j < data.length; j++) {
          data[j] ^= moves[i].data[j];
        }
      }

      return onRandomize(state.state, data);
    } else if (state.type_ == DType.CommittingSecret) {
      data = new bytes(state.statuses.length + moves.length * 32);

      for (i = 0; i < moves.length; i++) {
        data[moves[i].player] = byte(i);

        for (j = 0; j < 32; j++) {
          data[state.statuses.length + i * 32 + j] = moves[i].data[j];
        }
      }

      return DState(DType.RevealingRandom, data, state.statuses, state.state);
    } else if (state.type_ == DType.RevealingSecret) {
      return onExchange(state.state, moves);
    }
  }
}
