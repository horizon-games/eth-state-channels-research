pragma solidity ^0.4.19;
pragma experimental ABIEncoderV2;

import './Arcadeum.sol';

contract DGame {
  uint constant META_STATE_DATA_LENGTH = 3;
  uint constant STATE_DATA_LENGTH = 1;

  string constant ETH_SIGN_PREFIX = '\x19Ethereum Signed Message:\n';
  string constant MESSAGE_LENGTH = '158'; // INVITATION.length + GAME_PREFIX.length + 40 + MATCH_PREFIX.length + 8 + SUBKEY_PREFIX.length + 40
  string constant INVITATION = 'Sign to play! This won\'t cost anything.\n';
  string constant GAME_PREFIX = '\nGame: 0x';
  string constant MATCH_PREFIX = '\nMatch: 0x';
  string constant SUBKEY_PREFIX = '\nPlayer: 0x';

  struct Match {
    DGame game;
    uint32 matchID;
    Player[2] players;
    Signature signature;
  }

  struct Player {
    address account;
    address subkey;
    Signature0 subkeySignature;
    bytes publicSeed;
  }

  struct MetaState {
    uint32 nonce;
    MetaTag tag;
    bytes32[META_STATE_DATA_LENGTH] data;
    State state;
  }

  enum MetaTag {
    Playing,
    CommittingRandomness,
    RevealingRandomness,
    CommittingSecret,
    RevealingSecret
  }

  struct MetaMove {
    Move move;
    Signature signature;
  }

  struct Signature {
    uint8 v;
    bytes32 r;
    bytes32 s;
  }

  // XXX: https://github.com/ethereum/solidity/issues/3275
  struct Signature0 {
    uint8 v;
    bytes32 r;
    bytes32 s;
  }

  struct State {
    uint32 tag;
    bytes32[STATE_DATA_LENGTH] data;
  }

  struct Move {
    uint8 playerID;
    bytes data;
  }

  function DGame(Arcadeum arcadeum) internal {
    owner = arcadeum;
  }

  function secretSeedRating(address account, bytes secretSeed) public pure returns (uint32) {
    return secretSeedRatingInternal(account, secretSeed);
  }

  function publicSeed(address account, bytes secretSeed) public pure returns (bytes) {
    return publicSeedInternal(account, secretSeed);
  }

  function isSubkeySigned(Match dMatch, uint playerID) public pure returns (bool) {
    bytes memory gameString;
    bytes memory matchString;
    bytes memory subkeyString;
    uint i;
    uint8 b;
    uint8 hi;
    uint8 lo;

    gameString = new bytes(40);
    matchString = new bytes(8);
    subkeyString = new bytes(40);

    for (i = 0; i < 20; i++) {
      b = uint8(bytes20(address(dMatch.game))[i]);
      hi = b / 16;
      lo = b % 16;

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

      gameString[2 * i] = byte(hi);
      gameString[2 * i + 1] = byte(lo);

      b = uint8(bytes20(address(dMatch.players[playerID].subkey))[i]);
      hi = b / 16;
      lo = b % 16;

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

      subkeyString[2 * i] = byte(hi);
      subkeyString[2 * i + 1] = byte(lo);
    }

    for (i = 0; i < 4; i++) {
      b = uint8(bytes4(dMatch.matchID)[i]);
      hi = b / 16;
      lo = b % 16;

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

      matchString[2 * i] = byte(hi);
      matchString[2 * i + 1] = byte(lo);
    }

    return ecrecover(keccak256(ETH_SIGN_PREFIX, MESSAGE_LENGTH, INVITATION, GAME_PREFIX, gameString, MATCH_PREFIX, matchString, SUBKEY_PREFIX, subkeyString), dMatch.players[playerID].subkeySignature.v, dMatch.players[playerID].subkeySignature.r, dMatch.players[playerID].subkeySignature.s) == dMatch.players[playerID].account;
  }

  function initialState(Match dMatch) public pure returns (MetaState) {
    return initialStateInternal(dMatch);
  }

  function winner(MetaState mState) public pure returns (uint) {
    if (mState.tag == MetaTag.Playing) {
      return winnerInternal(mState.state);

    } else if (mState.tag == MetaTag.CommittingRandomness) {
      return 0;

    } else if (mState.tag == MetaTag.RevealingRandomness) {
      return 0;

    } else if (mState.tag == MetaTag.CommittingSecret) {
      return 0;

    } else if (mState.tag == MetaTag.RevealingSecret) {
      return 0;

    }
  }

  function nextPlayers(MetaState mState) public pure returns (uint) {
    if (winner(mState) != 0) {
      return 0;
    }

    if (mState.tag == MetaTag.Playing) {
      return nextPlayersInternal(mState.state);

    } else if (mState.tag == MetaTag.CommittingRandomness) {
      return 3;

    } else if (mState.tag == MetaTag.RevealingRandomness) {
      return 3;

    } else if (mState.tag == MetaTag.CommittingSecret) {
      return 3;

    } else if (mState.tag == MetaTag.RevealingSecret) {
      return 3;

    }
  }

  function isMoveLegal(MetaState mState, Move move) public pure returns (bool) {
    uint next;

    next = nextPlayers(mState);

    if (next != 3 && next != 1 + move.playerID) {
      return false;
    }

    if (mState.tag == MetaTag.Playing) {
      return isMoveLegalInternal(mState.state, move);

    } else if (mState.tag == MetaTag.CommittingRandomness) {
      return move.data.length == 32;

    } else if (mState.tag == MetaTag.RevealingRandomness) {
      if (move.data.length != uint(mState.data[0])) {
        return false;
      }

      return keccak256(move.data) == mState.data[1 + move.playerID];

    } else if (mState.tag == MetaTag.CommittingSecret) {
      return move.data.length == 32;

    } else if (mState.tag == MetaTag.RevealingSecret) {
      return keccak256(move.data) == mState.data[move.playerID];

    }
  }

  // XXX: https://github.com/ethereum/solidity/issues/3516
  function nextState1XXX(MetaState mState, Move move) public pure returns (uint32, MetaTag, bytes32[META_STATE_DATA_LENGTH], uint32, bytes32[STATE_DATA_LENGTH]) {
    MetaState memory next;

    next = nextState1(mState, move);

    return (mState.nonce, mState.tag, mState.data, mState.state.tag, mState.state.data);
  }

  // XXX: https://github.com/ethereum/solidity/issues/3516
  function nextState2XXX(MetaState mState, Move moveA, Move moveB) public pure returns (uint32, MetaTag, bytes32[META_STATE_DATA_LENGTH], uint32, bytes32[STATE_DATA_LENGTH]) {
    MetaState memory next;

    next = nextState2(mState, moveA, moveB);

    return (mState.nonce, mState.tag, mState.data, mState.state.tag, mState.state.data);
  }

  // XXX: https://github.com/ethers-io/ethers.js/issues/119
  function nextState1(MetaState mState, Move move) public pure returns (MetaState) {
    MetaState memory next;

    if (mState.tag == MetaTag.Playing) {
      next = nextStateInternal(mState.state, move);

    } else if (mState.tag == MetaTag.CommittingRandomness) {
      assert(false);

    } else if (mState.tag == MetaTag.RevealingRandomness) {
      assert(false);

    } else if (mState.tag == MetaTag.CommittingSecret) {
      assert(false);

    } else if (mState.tag == MetaTag.RevealingSecret) {
      assert(false);

    }

    next.nonce = mState.nonce + 1;

    return next;
  }

  // XXX: https://github.com/ethers-io/ethers.js/issues/119
  function nextState2(MetaState mState, Move moveA, Move moveB) public pure returns (MetaState) {
    MetaState memory next;
    Move[2] memory moves;
    bytes memory data;
    uint i;

    require(moveA.playerID != moveB.playerID);

    moves[moveA.playerID] = moveA;
    moves[moveB.playerID] = moveB;

    if (mState.tag == MetaTag.Playing) {
      next = nextStateInternal(mState.state, moves[0], moves[1]);

    } else if (mState.tag == MetaTag.CommittingRandomness) {
      for (i = 0; i < 32; i++) {
        if (moves[0].data[i] != moves[1].data[i]) {
          break;
        }
      }

      if (i == 32) {
        next = mState;

      } else {
        next.tag = MetaTag.RevealingRandomness;
        next.data[0] = mState.data[0];

        for (i = 0; i < 32; i++) {
          next.data[1] |= bytes32(moves[0].data[i]) << ((31 - i) * 8);
          next.data[2] |= bytes32(moves[1].data[i]) << ((31 - i) * 8);
        }

        next.state = mState.state;

      }

    } else if (mState.tag == MetaTag.RevealingRandomness) {
      data = new bytes(uint(mState.data[0]));

      for (i = 0; i < data.length; i++) {
        data[i] = moves[0].data[i] ^ moves[1].data[i];
      }

      next = onRandomizeInternal(data, mState.state);

    } else if (mState.tag == MetaTag.CommittingSecret) {
      for (i = 0; i < 32; i++) {
        if (moves[0].data[i] != moves[1].data[i]) {
          break;
        }
      }

      if (i == 32) {
        next = mState;

      } else {
        next.tag = MetaTag.RevealingSecret;

        for (i = 0; i < 32; i++) {
          next.data[0] |= bytes32(moves[0].data[i]) << ((31 - i) * 8);
          next.data[1] |= bytes32(moves[1].data[i]) << ((31 - i) * 8);
        }

        next.state = mState.state;

      }

    } else if (mState.tag == MetaTag.RevealingSecret) {
      next = onRevealInternal(moves[0].data, moves[1].data, mState.state);

    }

    next.nonce = mState.nonce + 1;

    return next;
  }

  modifier restricted() { require(msg.sender == address(owner)); _; }

  function playerWonRestricted(Match dMatch, uint winnerID) public restricted {
    onPlayerWonInternal(dMatch, winnerID);
  }

  function playerCheatedRestricted(Match dMatch, uint cheaterID) public restricted {
    onPlayerCheatedInternal(dMatch, cheaterID);
  }

  function secretSeedRatingInternal(address, bytes) internal pure returns (uint32) {
    return 0;
  }

  function publicSeedInternal(address, bytes secretSeed) internal pure returns (bytes) {
    return secretSeed;
  }

  function initialStateInternal(Match) internal pure returns (MetaState) {
    State memory state;

    return meta(state);
  }

  function winnerInternal(State state) internal pure returns (uint);
  function nextPlayersInternal(State state) internal pure returns (uint);
  function isMoveLegalInternal(State state, Move move) internal pure returns (bool);

  function nextStateInternal(State, Move) internal pure returns (MetaState) {
    assert(false);
  }

  function nextStateInternal(State, Move, Move) internal pure returns (MetaState) {
    assert(false);
  }

  function onRandomizeInternal(bytes, State) internal pure returns (MetaState) {
    assert(false);
  }

  function onRevealInternal(bytes, bytes, State) internal pure returns (MetaState) {
    assert(false);
  }

  function onPlayerWonInternal(Match, uint) internal {
  }

  function onPlayerCheatedInternal(Match, uint) internal {
  }

  function meta(State state) internal pure returns (MetaState) {
    MetaState memory mState;

    mState.tag = MetaTag.Playing;
    mState.state = state;

    return mState;
  }

  function randomize(uint nBytes, State state) internal pure returns (MetaState) {
    MetaState memory mState;

    require(nBytes < 256);

    mState.tag = MetaTag.CommittingRandomness;
    mState.data[0] = bytes32(nBytes);
    mState.state = state;

    return mState;
  }

  function commit(State state) internal pure returns (MetaState) {
    MetaState memory mState;

    mState.tag = MetaTag.CommittingSecret;
    mState.state = state;

    return mState;
  }

  Arcadeum owner;
}
