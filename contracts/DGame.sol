pragma solidity ^0.4.19;
pragma experimental ABIEncoderV2;

contract DGame {
  struct Match {
    DGame game;
    uint32 matchID;
    Player[] players;
    Signature signature; // matchmaker.sign(hash(game, matchID, accounts))
  }

  struct Player {
    address account;
    address subkey;
    Signature1 subkeySignature; // account.sign(hash(format(game, matchID, subkey)))
    bytes seed;
    Signature2 seedSignature; // subkey.sign(hash(seed))
  }

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

  struct MetaMove {
    bytes32 mStateHash; // hash(mNonce, mTag, mData, mStatuses, tag, data, statuses)
    Move move;
    Signature signature; // subkey.sign(hash(mStateHash, playerID, data))
  }

  struct Signature {
    uint8 v;
    bytes32 r;
    bytes32 s;
  }

  // https://github.com/ethereum/solidity/issues/3275
  struct Signature1 {
    uint8 v;
    bytes32 r;
    bytes32 s;
  }

  // https://github.com/ethereum/solidity/issues/3275
  struct Signature2 {
    uint8 v;
    bytes32 r;
    bytes32 s;
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

  function matchmaker(Match match_) public pure returns (address) {
    address[] memory players;
    uint i;

    players = new address[](match_.players.length);

    for (i = 0; i < players.length; i++) {
      players[i] = match_.players[i].account;
    }

    return ecrecover(keccak256(match_.game, match_.matchID, players), match_.signature.v, match_.signature.r, match_.signature.s);
  }

  function playerRank(address account, bytes seed) public pure returns (uint32) {
    return playerRankInternal(account, seed);
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

  function isMoveSigned(Match match_, MetaState mState, MetaMove mMove) public pure returns (bool) {
    bytes memory message;
    byte b;
    uint8 hi;
    uint8 lo;
    uint i;

    message = "Sign to play!\n\nGame: 0x????????????????????????????????????????\nMatch: 0x????????\nSubkey: 0x????????????????????????????????????????\n";

    for (i = 0; i < 20; i++) {
      b = bytes20(address(match_.game))[i];
      hi = uint8(b) / 16;
      lo = uint8(b) % 16;

      if (hi < 10) {
        hi += 48;
      } else {
        hi += 87;
      }

      if (lo < 10) {
        lo += 48;
      } else {
        lo += 87;
      }

      message[23 + i] = byte(hi);
      message[24 + i] = byte(lo);

      b = bytes20(match_.players[mMove.move.playerID].subkey)[i];
      hi = uint8(b) / 16;
      lo = uint8(b) % 16;

      if (hi < 10) {
        hi += 48;
      } else {
        hi += 87;
      }

      if (lo < 10) {
        lo += 48;
      } else {
        lo += 87;
      }

      message[92 + i] = byte(hi);
      message[93 + i] = byte(lo);
    }

    for (i = 0; i < 4; i++) {
      hi = uint8(b) / 16;
      lo = uint8(b) % 16;

      if (hi < 10) {
        hi += 48;
      } else {
        hi += 87;
      }

      if (lo < 10) {
        lo += 48;
      } else {
        lo += 87;
      }

      message[73 + i] = byte(hi);
      message[73 + i] = byte(lo);
    }

    if (ecrecover(keccak256(message), match_.players[mMove.move.playerID].subkeySignature.v, match_.players[mMove.move.playerID].subkeySignature.r, match_.players[mMove.move.playerID].subkeySignature.s) != match_.players[mMove.move.playerID].account) {
      return false;
    }

    if (keccak256(mState.nonce, mState.tag, mState.data, mState.statuses, mState.state.tag, mState.state.data, mState.state.statuses) != mMove.mStateHash) {
      return false;
    }

    if (ecrecover(keccak256(mMove.mStateHash, mMove.move.playerID, mMove.move.data), mMove.signature.v, mMove.signature.r, mMove.signature.s) != match_.players[mMove.move.playerID].subkey) {
      return false;
    }

    return true;
  }

  function isMoveLegal(MetaState mState, Move move) public pure returns (bool) {
    bytes32 hash;
    uint i;
    uint j;

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
      i = uint(mState.data[1 + move.playerID]);

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
      i = uint(mState.data[move.playerID]);

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
    uint i;
    uint j;

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

  function playerRankInternal(address, bytes) internal pure returns (uint32) {
    return 0;
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
    uint i;

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
    uint i;

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
