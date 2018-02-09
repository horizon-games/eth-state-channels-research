pragma solidity ^0.4.19;
pragma experimental ABIEncoderV2;

contract DGame {
  struct MetaState {
    uint16 nonce;
    MetaType tag;
    bytes data;
    Status[] statuses;
    State state;
  }

  enum MetaType {
    None,
    CommittingRandom,
    RevealingRandom,
    CommittingSecret,
    RevealingSecret
  }

  struct State {
    uint8 tag;
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
    uint8 playerID;
    bytes data;
  }

  struct Signature {
    uint8 v;
    bytes32 r;
    bytes32 s;
  }

  function playerStatus(MetaState mState, uint8 playerID) public pure returns (Status) {
    Status status;

    status = Status.Playing;

    if (mState.statuses.length > playerID) {
      status = mState.statuses[playerID];
    }

    if (status == Status.Playing) {
      if (mState.state.statuses.length > playerID) {
        status = mState.state.statuses[playerID];
      }

      if (status == Status.Playing) {
        status = playerStatusInternal(mState.state, playerID);
      }
    }

    return status;
  }

  function isMoveLegal(MetaState mState, Move move) public pure returns (bool) {
    bytes32 hash;
    uint8 i;
    uint8 j;

    if (playerStatus(mState, move.playerID) != Status.Moving) {
      return false;
    }

    if (mState.tag == MetaType.None) {
      return isMoveLegalInternal(mState.state, move);

    } else if (mState.tag == MetaType.CommittingRandom) {
      return move.data.length == 32;

    } else if (mState.tag == MetaType.RevealingRandom) {
      if (move.data.length != uint(mState.data[0])) {
        return false;
      }

      hash = keccak256(move.data);
      i = uint8(mState.data[1 + move.playerID]);

      for (j = 0; j < 32; j++) {
        if (hash[j] != mState.data[1 + mState.statuses.length + i * 32 + j]) {
          return false;
        }
      }

      return true;

    } else if (mState.tag == MetaType.CommittingSecret) {
      return move.data.length == 32;

    } else if (mState.tag == MetaType.RevealingSecret) {
      hash = keccak256(move.data);
      i = uint8(mState.data[move.playerID]);

      for (j = 0; j < 32; j++) {
        if (hash[j] != mState.data[mState.statuses.length + i * 32 + j]) {
          return false;
        }
      }

      return true;
    }
  }

  function nextState(MetaState mState, Move[] moves) public pure returns (MetaState) {
    MetaState memory next;
    bytes memory data;
    uint8 i;
    uint8 j;

    if (mState.tag == MetaType.None) {
      next = nextStateInternal(mState.state, moves);

    } else if (mState.tag == MetaType.CommittingRandom) {
      data = new bytes(1 + mState.statuses.length + moves.length * 32);
      data[0] = mState.data[0];

      for (i = 0; i < moves.length; i++) {
        data[1 + moves[i].playerID] = byte(i);

        for (j = 0; j < 32; j++) {
          data[1 + mState.statuses.length + i * 32 + j] = moves[i].data[j];
        }
      }

      next = MetaState(0, MetaType.RevealingRandom, data, mState.statuses, mState.state);

    } else if (mState.tag == MetaType.RevealingRandom) {
      data = new bytes(uint(mState.data[0]));

      for (i = 0; i < moves.length; i++) {
        for (j = 0; j < data.length; j++) {
          data[j] ^= moves[i].data[j];
        }
      }

      next = onRandomizeInternal(mState.state, data);

    } else if (mState.tag == MetaType.CommittingSecret) {
      data = new bytes(mState.statuses.length + moves.length * 32);

      for (i = 0; i < moves.length; i++) {
        data[moves[i].playerID] = byte(i);

        for (j = 0; j < 32; j++) {
          data[mState.statuses.length + i * 32 + j] = moves[i].data[j];
        }
      }

      next = MetaState(0, MetaType.RevealingSecret, data, mState.statuses, mState.state);

    } else if (mState.tag == MetaType.RevealingSecret) {
      next = onExchangeInternal(mState.state, moves);
    }

    next.nonce = mState.nonce + 1;
    return next;
  }

  function playerStatusInternal(State state, uint8 playerID) internal pure returns (Status);
  function isMoveLegalInternal(State state, Move move) internal pure returns (bool);
  function nextStateInternal(State state, Move[] moves) internal pure returns (MetaState);

  function onRandomizeInternal(State state, bytes) internal pure returns (MetaState) {
    return id(state);
  }

  function onExchangeInternal(State state, Move[]) internal pure returns (MetaState) {
    return id(state);
  }

  function id(State state) internal pure returns (MetaState) {
    return MetaState(0, MetaType.None, new bytes(0), new Status[](0), state);
  }

  function randomize(State state, uint8 numBytes, uint8[] playerIDs) internal pure returns (MetaState) {
    bytes memory data;
    Status[] memory statuses;
    Status status;
    uint8 i;

    data = new bytes(1);
    data[0] = byte(numBytes);

    statuses = new Status[](state.statuses.length);

    for (i = 0; i < statuses.length; i++) {
      status = state.statuses[i];

      if (status == Status.Playing || status == Status.Moving) {
        status = Status.Waiting;
      }

      statuses[i] = status;
    }

    for (i = 0; i < playerIDs.length; i++) {
      status = statuses[playerIDs[i]];

      if (status == Status.Waiting) {
        statuses[playerIDs[i]] = Status.Moving;
      }
    }

    return MetaState(0, MetaType.CommittingRandom, data, statuses, state);
  }

  function exchange(State state, uint8[] playerIDs) internal pure returns (MetaState) {
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

    for (i = 0; i < playerIDs.length; i++) {
      status = statuses[playerIDs[i]];

      if (status == Status.Waiting) {
        statuses[playerIDs[i]] = Status.Moving;
      }
    }

    return MetaState(0, MetaType.CommittingSecret, new bytes(0), statuses, state);
  }
}
