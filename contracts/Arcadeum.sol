pragma solidity ^0.4.19;
pragma experimental ABIEncoderV2;

import './DGame.sol';

contract Arcadeum {
  uint private constant WITHDRAWAL_TIME = 10 minutes;
  uint private constant STOP_WITHDRAWAL_GAS = 21000; // XXX

  // *** THIS MUST MATCH subkeyMessage ***
  string private constant ETH_SIGN_PREFIX = '\x19Ethereum Signed Message:\n';
  string private constant MESSAGE_LENGTH = '91'; // 40 + 11 + 40
  string private constant MESSAGE_PREFIX = 'Sign to play! This won\'t cost anything.\n';
  string private constant PLAYER_PREFIX = '\nPlayer: 0x';

  struct Match {
    DGame game;
    uint32 matchID;
    uint timestamp;
    uint8 playerID;
    Player[2] players;
    Signature matchSignature;
    // XXX: https://github.com/ethereum/solidity/issues/3275#issuecomment-365087323
    SubkeySignature opponentSubkeySignature;
    TimestampSignature opponentTimestampSignature;
  }

  struct Player {
    uint32 seedRating;
    bytes publicSeed;
  }

  struct Move {
    DGame.Move move;
    Signature signature;
  }

  struct Signature {
    uint8 v;
    bytes32 r;
    bytes32 s;
  }

  // XXX: https://github.com/ethereum/solidity/issues/3275#issuecomment-365087323
  struct SubkeySignature {
    uint8 v;
    bytes32 r;
    bytes32 s;
  }

  // XXX: https://github.com/ethereum/solidity/issues/3275#issuecomment-365087323
  struct TimestampSignature {
    uint8 v;
    bytes32 r;
    bytes32 s;
  }

  function Arcadeum() public {
    owner = msg.sender;
  }

  mapping(address => uint) public balance;

  function deposit() external payable {
    balance[msg.sender] += msg.value;

    balanceChanged(msg.sender);
  }

  function isWithdrawing(address account) public view returns (bool) {
    return withdrawalTime[account] != 0;
  }

  function startWithdrawal() external {
    require(balance[msg.sender] > 0);

    withdrawalTime[msg.sender] = now;

    withdrawalStarted(msg.sender);
  }

  function canFinishWithdrawal(address account) public view returns (bool) {
    uint startTime;

    startTime = withdrawalTime[account];

    if (startTime == 0) {
      return false;

    } else if (now < startTime + WITHDRAWAL_TIME) {
      return false;
    }

    return true;
  }

  function finishWithdrawal() external {
    uint value;

    require(canFinishWithdrawal(msg.sender));

    delete withdrawalTime[msg.sender];

    value = balance[msg.sender];
    delete balance[msg.sender];
    msg.sender.transfer(value);

    balanceChanged(msg.sender);
  }

  // XXX: https://github.com/ethereum/solidity/issues/3275#issuecomment-365087323
  function couldStopWithdrawal(DGame game, uint32 matchID, uint timestamp, TimestampSignature, SubkeySignature) public view returns (bool) {
    uint expiryTime;
    bytes24 gameMatchID;

    expiryTime = timestamp + game.matchDuration();

    if (now >= expiryTime) {
      return false;
    }

    gameMatchID = (bytes24(address(game)) << 32) | bytes24(matchID);

    if (isMatchFinal[gameMatchID]) {
      return false;
    }

    return true;
  }

  // XXX: https://github.com/ethereum/solidity/issues/3275#issuecomment-365087323
  function canStopWithdrawal(DGame game, uint32 matchID, uint timestamp, TimestampSignature timestampSignature, SubkeySignature subkeySignature) public view returns (bool) {
    address account;

    account = playerAccount(game, matchID, timestamp, timestampSignature, subkeySignature);

    return isWithdrawing(account) && couldStopWithdrawal(game, matchID, timestamp, timestampSignature, subkeySignature);
  }

  // XXX: https://github.com/ethereum/solidity/issues/3275#issuecomment-365087323
  function stopWithdrawal(DGame game, uint32 matchID, uint timestamp, TimestampSignature timestampSignature, SubkeySignature subkeySignature) external {
    address account;
    uint value;

    require(canStopWithdrawal(game, matchID, timestamp, timestampSignature, subkeySignature));

    account = playerAccount(game, matchID, timestamp, timestampSignature, subkeySignature);
    delete withdrawalTime[account];
    value = STOP_WITHDRAWAL_GAS * tx.gasprice;

    if (value > balance[account]) {
      value = balance[account];
    }

    balance[account] -= value;
    balance[owner] += value;

    balanceChanged(account);
    balanceChanged(owner);
    withdrawalStopped(account);
  }

  event balanceChanged(address indexed account);
  event withdrawalStarted(address indexed account);
  event withdrawalStopped(address indexed account);

  function canClaimReward(Match aMatch, DGame.MetaState metaState, Move loserMove, DGame.Move[] winnerMoves) public view returns (bool) {
    bytes24 gameMatchID;
    address opponent;
    bool isLegal;
    DGame.Winner winner;
    DGame.NextPlayers nextPlayers;
    DGame.MetaState memory nextState;
    uint i;

    gameMatchID = (bytes24(address(aMatch.game)) << 32) | bytes24(aMatch.matchID);

    if (isMatchFinal[gameMatchID]) {
      return false;
    }

    if (matchMaker(aMatch, msg.sender) != owner) {
      return false;
    }

    opponent = playerAccount(aMatch.game, aMatch.matchID, aMatch.timestamp, aMatch.opponentTimestampSignature, aMatch.opponentSubkeySignature);

    if (moveMaker(metaState, loserMove, aMatch.opponentSubkeySignature) != opponent) {
      return false;
    }

    winner = aMatch.game.winner(metaState);

    if (winner == DGame.Winner.NONE) {
      (isLegal,) = aMatch.game.isMoveLegal(metaState, loserMove.move);

      if (!isLegal) {
        return false;
      }

      nextPlayers = aMatch.game.nextPlayers(metaState);

      if (nextPlayers != DGame.NextPlayers.BOTH) {
        nextState = aMatch.game.nextState(metaState, loserMove.move);
        i = 0;

      } else /* nextPlayers == DGame.NextPlayers.PLAYER_0 || nextPlayers == DGame.NextPlayers.PLAYER_1 */ {
        if (winnerMoves[0].playerID != aMatch.playerID) {
          return false;
        }

        (isLegal,) = aMatch.game.isMoveLegal(metaState, winnerMoves[0]);

        if (!isLegal) {
          return false;
        }

        nextState = aMatch.game.nextState(metaState, loserMove.move, winnerMoves[0]);
        i = 1;
      }
    }

    for (; winner == DGame.Winner.NONE; i++) {
      if (winnerMoves[i].playerID != aMatch.playerID) {
        return false;
      }

      nextPlayers = aMatch.game.nextPlayers(nextState);

      if (nextPlayers == DGame.NextPlayers.BOTH) {
        return false;
      }

      (isLegal,) = aMatch.game.isMoveLegal(nextState, winnerMoves[i]);

      if (!isLegal) {
        return false;
      }

      nextState = aMatch.game.nextState(nextState, winnerMoves[i]);
      winner = aMatch.game.winner(nextState);
    }

    if (winner == DGame.Winner.PLAYER_0) {
      if (aMatch.playerID != 0) {
        return false;
      }

    } else /* winner == DGame.Winner.PLAYER_1 */ {
      if (aMatch.playerID != 1) {
        return false;
      }
    }

    return true;
  }

  function claimReward(Match aMatch, DGame.MetaState metaState, Move loserMove, DGame.Move[] winnerMoves) external {
    bytes24 gameMatchID;
    uint32 winnerSeedRating;
    uint32 opponentSeedRating;

    require(canClaimReward(aMatch, metaState, loserMove, winnerMoves));

    gameMatchID = (bytes24(address(aMatch.game)) << 32) | bytes24(aMatch.matchID);
    isMatchFinal[gameMatchID] = true;

    winnerSeedRating = aMatch.players[aMatch.playerID].seedRating;
    opponentSeedRating = aMatch.players[1 - aMatch.playerID].seedRating;
    aMatch.game.registerWin(msg.sender, winnerSeedRating, opponentSeedRating, metaState.state, aMatch.playerID);

    rewardClaimed(aMatch.game, aMatch.matchID, msg.sender);
  }

  function canReportCheater(Match aMatch, DGame.MetaState metaState, Move cheaterMove) public view returns (bool) {
    bytes24 gameMatchID;
    address opponent;
    bool isLegal;

    gameMatchID = (bytes24(address(aMatch.game)) << 32) | bytes24(aMatch.matchID);

    if (isMatchFinal[gameMatchID]) {
      return false;
    }

    if (matchMaker(aMatch, msg.sender) != owner) {
      return false;
    }

    opponent = playerAccount(aMatch.game, aMatch.matchID, aMatch.timestamp, aMatch.opponentTimestampSignature, aMatch.opponentSubkeySignature);

    if (moveMaker(metaState, cheaterMove, aMatch.opponentSubkeySignature) != opponent) {
      return false;
    }

    (isLegal,) = aMatch.game.isMoveLegal(metaState, cheaterMove.move);

    if (isLegal) {
      return false;
    }

    return true;
  }

  function reportCheater(Match aMatch, DGame.MetaState metaState, Move cheaterMove) external {
    bytes24 gameMatchID;
    address opponent;
    uint value;

    require(canReportCheater(aMatch, metaState, cheaterMove));

    gameMatchID = (bytes24(address(aMatch.game)) << 32) | bytes24(aMatch.matchID);
    isMatchFinal[gameMatchID] = true;

    opponent = playerAccount(aMatch.game, aMatch.matchID, aMatch.timestamp, aMatch.opponentTimestampSignature, aMatch.opponentSubkeySignature);
    value = balance[opponent];
    delete balance[opponent];
    balance[msg.sender] += value / 2;
    balance[owner] += (value + 1) / 2;

    aMatch.game.registerCheat(opponent);

    balanceChanged(opponent);
    balanceChanged(msg.sender);
    balanceChanged(owner);
    cheaterReported(aMatch.game, aMatch.matchID, opponent);
  }

  event rewardClaimed(address indexed game, uint32 indexed matchID, address indexed account);
  event cheaterReported(address indexed game, uint32 indexed matchID, address indexed account);

  function subkeyMessage(address subkey) public pure returns (string) {
    bytes memory message;
    bytes20 subkeyBytes;
    uint i;
    uint8 b;
    uint8 hi;
    uint8 lo;

    message = new bytes(91);

    for (i = 0; i < 40; i++) {
      message[i] = bytes(MESSAGE_PREFIX)[i];
    }

    for (i = 0; i < 11; i++) {
      message[40 + i] = bytes(PLAYER_PREFIX)[i];
    }

    subkeyBytes = bytes20(subkey);

    for (i = 0; i < 20; i++) {
      b = uint8(subkeyBytes[i]);
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

      message[51 + 2 * i + 0] = byte(hi);
      message[51 + 2 * i + 1] = byte(lo);
    }

    return string(message);
  }

  // XXX: https://github.com/ethereum/solidity/issues/3275#issuecomment-365087323
  function subkeyParent(address subkey, SubkeySignature subkeySignature) public pure returns (address) {
    bytes20 subkeyBytes;
    bytes memory subkeyHex;
    uint i;
    uint8 b;
    uint8 hi;
    uint8 lo;
    bytes32 hash;

    subkeyBytes = bytes20(subkey);
    subkeyHex = new bytes(40);

    for (i = 0; i < 20; i++) {
      b = uint8(subkeyBytes[i]);
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

      subkeyHex[2 * i + 0] = byte(hi);
      subkeyHex[2 * i + 1] = byte(lo);
    }

    hash = keccak256(ETH_SIGN_PREFIX, MESSAGE_LENGTH, MESSAGE_PREFIX, PLAYER_PREFIX, subkeyHex);

    return ecrecover(hash, subkeySignature.v, subkeySignature.r, subkeySignature.s);
  }

  // XXX: https://github.com/ethereum/solidity/issues/3275#issuecomment-365087323
  function playerAccount(DGame game, uint32 matchID, uint timestamp, TimestampSignature timestampSignature, SubkeySignature subkeySignature) public pure returns (address) {
    bytes32 hash;
    address subkey;

    hash = keccak256(game, matchID, timestamp);
    subkey = ecrecover(hash, timestampSignature.v, timestampSignature.r, timestampSignature.s);

    return subkeyParent(subkey, subkeySignature);
  }

  function matchHash(DGame game, uint32 matchID, uint timestamp, address[2] accounts, uint32[2] seedRatings, bytes[2] publicSeeds) public pure returns (bytes32) {
    return keccak256(game, matchID, timestamp, accounts[0], accounts[1], seedRatings[0], seedRatings[1], publicSeeds[0], publicSeeds[1]);
  }

  function matchMaker(Match aMatch, address sender) private pure returns (address) {
    address[2] memory accounts;
    uint32[2] memory seedRatings;
    bytes[2] memory publicSeeds;
    bytes32 hash;

    accounts[aMatch.playerID] = sender;
    accounts[1 - aMatch.playerID] = playerAccount(aMatch.game, aMatch.matchID, aMatch.timestamp, aMatch.opponentTimestampSignature, aMatch.opponentSubkeySignature);
    seedRatings[0] = aMatch.players[0].seedRating;
    seedRatings[1] = aMatch.players[1].seedRating;
    publicSeeds[0] = aMatch.players[0].publicSeed;
    publicSeeds[1] = aMatch.players[1].publicSeed;
    hash = matchHash(aMatch.game, aMatch.matchID, aMatch.timestamp, accounts, seedRatings, publicSeeds);

    return ecrecover(hash, aMatch.matchSignature.v, aMatch.matchSignature.r, aMatch.matchSignature.s);
  }

  function stateHash(DGame.MetaState metaState) public pure returns (bytes32) {
    return keccak256(metaState.nonce, metaState.tag, metaState.data, metaState.state.tag, metaState.state.data);
  }

  // XXX: https://github.com/ethereum/solidity/issues/3275#issuecomment-365087323
  function moveMaker(DGame.MetaState metaState, Move move, SubkeySignature subkeySignature) public pure returns (address) {
    bytes32 hash;

    hash = keccak256(stateHash(metaState), move.move.playerID, move.move.data);

    return subkeyParent(ecrecover(hash, move.signature.v, move.signature.r, move.signature.s), subkeySignature);
  }

  address private owner;
  mapping(address => uint) private withdrawalTime;
  mapping(bytes24 => bool) private isMatchFinal;
}
