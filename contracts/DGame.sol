pragma solidity ^0.4.19;
pragma experimental ABIEncoderV2;

contract DGame {
  struct DState {
    uint16 nonce;
    DType type_;
    bytes data;
    Status[] statuses;
    State state;
  }

  enum DType {
    None,
    CommittingRandom,
    RevealingRandom,
    CommittingSecret,
    RevealingSecret
  }

  struct State {
    uint8 type_;
    bytes data;
    Status[] statuses;
  }

  enum Status {
    Playing,
    Waiting,
    Moving,
    Done,
    Won,
    Cheated
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

  function playerStatus(DState state, uint8 player) public pure returns (Status) {
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
        status = playerStatusInternal(state.state, player);
      }
    }

    return status;
  }

  function isMoveLegal(DState state, Move move) public pure returns (bool) {
    bytes32 hash;
    uint8 i;
    uint8 j;

    if (playerStatus(state, move.player) != Status.Moving) {
      return false;
    }

    if (state.type_ == DType.None) {
      return isMoveLegalInternal(state.state, move);
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

  function nextState(DState state, Move[] moves) public pure returns (DState) {
    DState memory next;
    bytes memory data;
    uint8 i;
    uint8 j;

    if (state.type_ == DType.None) {
      next = nextStateInternal(state.state, moves);
    } else if (state.type_ == DType.CommittingRandom) {
      data = new bytes(1 + state.statuses.length + moves.length * 32);
      data[0] = state.data[0];

      for (i = 0; i < moves.length; i++) {
        data[1 + moves[i].player] = byte(i);

        for (j = 0; j < 32; j++) {
          data[1 + state.statuses.length + i * 32 + j] = moves[i].data[j];
        }
      }

      next = DState(0, DType.RevealingRandom, data, state.statuses, state.state);
    } else if (state.type_ == DType.RevealingRandom) {
      data = new bytes(uint(state.data[0]));

      for (i = 0; i < moves.length; i++) {
        for (j = 0; j < data.length; j++) {
          data[j] ^= moves[i].data[j];
        }
      }

      next = onRandomizeInternal(state.state, data);
    } else if (state.type_ == DType.CommittingSecret) {
      data = new bytes(state.statuses.length + moves.length * 32);

      for (i = 0; i < moves.length; i++) {
        data[moves[i].player] = byte(i);

        for (j = 0; j < 32; j++) {
          data[state.statuses.length + i * 32 + j] = moves[i].data[j];
        }
      }

      next = DState(0, DType.RevealingSecret, data, state.statuses, state.state);
    } else if (state.type_ == DType.RevealingSecret) {
      next = onExchangeInternal(state.state, moves);
    }

    next.nonce = state.nonce + 1;
    return next;
  }

  function playerStatusInternal(State state, uint8 player) internal pure returns (Status);
  function isMoveLegalInternal(State state, Move move) internal pure returns (bool);
  function nextStateInternal(State state, Move[] moves) internal pure returns (DState);

  function onRandomizeInternal(State state, bytes) internal pure returns (DState) {
    return id(state);
  }

  function onExchangeInternal(State state, Move[]) internal pure returns (DState) {
    return id(state);
  }

  function id(State state) internal pure returns (DState) {
    return DState(0, DType.None, new bytes(0), new Status[](0), state);
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

    return DState(0, DType.CommittingRandom, data, statuses, state);
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

    return DState(0, DType.CommittingSecret, new bytes(0), statuses, state);
  }
}
