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

  // XXX: https://github.com/ethereum/solidity/issues/3270
  // *** THIS MUST MATCH DGame.sol ***
  uint private constant PUBLIC_SEED_LENGTH = 1;

  struct Match {
    DGame game;
    uint timestamp;
    uint8 playerID;
    Player[2] players;
    Signature matchSignature;
    // XXX: https://github.com/ethereum/solidity/issues/3275#issuecomment-365087323
    SubkeySignature opponentSubkeySignature;
  }

  struct Player {
    uint32 seedRating;
    // XXX: https://github.com/ethereum/solidity/issues/3270
    bytes32[PUBLIC_SEED_LENGTH] publicSeed;
    // XXX: https://github.com/ethereum/solidity/issues/3275#issuecomment-365087323
    TimestampSignature timestampSignature;
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
  mapping(address => uint) public withdrawalTime;

  function deposit() external payable {
    balance[msg.sender] += msg.value;

    emit balanceChanged(msg.sender);
  }

  function isWithdrawing(address account) public view returns (bool) {
    return withdrawalTime[account] != 0;
  }

  function startWithdrawal() external {
    require(balance[msg.sender] > 0);

    withdrawalTime[msg.sender] = now + WITHDRAWAL_TIME;

    emit withdrawalStarted(msg.sender);
  }

  function canFinishWithdrawal(address account) public view returns (bool) {
    uint time;

    time = withdrawalTime[account];

    if (time == 0) {
      return false;

    } else if (now < time) {
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

    emit balanceChanged(msg.sender);
  }

  // XXX: abigen: Failed to generate ABI binding: unsupported arg type: tuple
  function couldStopWithdrawalXXX(uint timestamp, uint8 timestampV, bytes32 timestampR, bytes32 timestampS, uint8 subkeyV, bytes32 subkeyR, bytes32 subkeyS) public view returns (bool) {
    return couldStopWithdrawal(timestamp, TimestampSignature(timestampV, timestampR, timestampS), SubkeySignature(subkeyV, subkeyR, subkeyS));
  }

  // XXX: https://github.com/ethereum/solidity/issues/3275#issuecomment-365087323
  function couldStopWithdrawal(uint timestamp, TimestampSignature timestampSignature, SubkeySignature) public view returns (bool) {
    if (now >= timestamp) {
      return false;
    }

    if (isTimestampInvalid(timestamp, timestampSignature)) {
      return false;
    }

    return true;
  }

  // XXX: abigen: Failed to generate ABI binding: unsupported arg type: tuple
  function canStopWithdrawalXXX(uint timestamp, uint8 timestampV, bytes32 timestampR, bytes32 timestampS, uint8 subkeyV, bytes32 subkeyR, bytes32 subkeyS) public view returns (bool) {
    return canStopWithdrawal(timestamp, TimestampSignature(timestampV, timestampR, timestampS), SubkeySignature(subkeyV, subkeyR, subkeyS));
  }

  // XXX: https://github.com/ethereum/solidity/issues/3275#issuecomment-365087323
  function canStopWithdrawal(uint timestamp, TimestampSignature timestampSignature, SubkeySignature subkeySignature) public view returns (bool) {
    address account;

    account = playerAccount(timestamp, timestampSignature, subkeySignature);

    return isWithdrawing(account) && couldStopWithdrawal(timestamp, timestampSignature, subkeySignature);
  }

  // XXX: abigen: Failed to generate ABI binding: unsupported arg type: tuple
  function stopWithdrawalXXX(uint timestamp, uint8 timestampV, bytes32 timestampR, bytes32 timestampS, uint8 subkeyV, bytes32 subkeyR, bytes32 subkeyS) public {
    stopWithdrawal(timestamp, TimestampSignature(timestampV, timestampR, timestampS), SubkeySignature(subkeyV, subkeyR, subkeyS));
  }

  // XXX: https://github.com/ethereum/solidity/issues/3199#issuecomment-365035663
  // XXX: https://github.com/ethereum/solidity/issues/3275#issuecomment-365087323
  function stopWithdrawal(uint timestamp, TimestampSignature timestampSignature, SubkeySignature subkeySignature) public {
    address account;
    uint value;

    require(canStopWithdrawal(timestamp, timestampSignature, subkeySignature));

    account = playerAccount(timestamp, timestampSignature, subkeySignature);
    delete withdrawalTime[account];
    value = STOP_WITHDRAWAL_GAS * tx.gasprice;

    if (value > balance[account]) {
      value = balance[account];
    }

    balance[account] -= value;
    balance[owner] += value;

    emit balanceChanged(account);
    emit balanceChanged(owner);
    emit withdrawalStopped(account);
  }

  event balanceChanged(address indexed account);
  event withdrawalStarted(address indexed account);
  event withdrawalStopped(address indexed account);

  function canClaimReward(Match aMatch, DGame.MetaState metaState, Move loserMove, DGame.Move[] winnerMoves) public view returns (bool) {
    address opponent;
    bool isLegal;
    DGame.Winner winner;
    DGame.NextPlayers nextPlayers;
    DGame.MetaState memory nextState;
    uint i;

    if (isTimestampInvalid(aMatch.timestamp, aMatch.players[aMatch.playerID].timestampSignature)) {
      return false;
    }

    if (matchMaker(aMatch, msg.sender) != owner) {
      return false;
    }

    opponent = playerAccount(aMatch.timestamp, aMatch.players[1 - aMatch.playerID].timestampSignature, aMatch.opponentSubkeySignature);

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
        // XXX: https://github.com/ethereum/solidity/issues/3516
        (nextState.nonce, nextState.tag, nextState.data, nextState.state.tag, nextState.state.data) = aMatch.game.nextStateXXX(metaState, loserMove.move);
        i = 0;

      } else /* nextPlayers == DGame.NextPlayers.PLAYER_0 || nextPlayers == DGame.NextPlayers.PLAYER_1 */ {
        if (winnerMoves[0].playerID != aMatch.playerID) {
          return false;
        }

        (isLegal,) = aMatch.game.isMoveLegal(metaState, winnerMoves[0]);

        if (!isLegal) {
          return false;
        }

        // XXX: https://github.com/ethereum/solidity/issues/3516
        (nextState.nonce, nextState.tag, nextState.data, nextState.state.tag, nextState.state.data) = aMatch.game.nextStateXXX(metaState, loserMove.move, winnerMoves[0]);
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

      // XXX: https://github.com/ethereum/solidity/issues/3516
      (nextState.nonce, nextState.tag, nextState.data, nextState.state.tag, nextState.state.data) = aMatch.game.nextStateXXX(nextState, winnerMoves[i]);
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

  // XXX: https://github.com/ethereum/solidity/issues/3199#issuecomment-365035663
  function claimReward(Match aMatch, DGame.MetaState metaState, Move loserMove, DGame.Move[] winnerMoves) public {
    uint32 winnerSeedRating;
    uint32 opponentSeedRating;

    require(canClaimReward(aMatch, metaState, loserMove, winnerMoves));

    invalidateTimestamp(aMatch.timestamp, aMatch.players[0].timestampSignature);
    invalidateTimestamp(aMatch.timestamp, aMatch.players[1].timestampSignature);

    winnerSeedRating = aMatch.players[aMatch.playerID].seedRating;
    opponentSeedRating = aMatch.players[1 - aMatch.playerID].seedRating;
    aMatch.game.registerWin(msg.sender, winnerSeedRating, opponentSeedRating, metaState.state, aMatch.playerID);

    emit rewardClaimed(msg.sender, timestampSubkey(aMatch.timestamp, aMatch.players[aMatch.playerID].timestampSignature), aMatch.timestamp);
  }

  function canReportCheater(Match aMatch, DGame.MetaState metaState, Move cheaterMove) public view returns (bool) {
    address opponent;
    bool isLegal;

    if (isTimestampInvalid(aMatch.timestamp, aMatch.players[aMatch.playerID].timestampSignature)) {
      return false;
    }

    if (matchMaker(aMatch, msg.sender) != owner) {
      return false;
    }

    opponent = playerAccount(aMatch.timestamp, aMatch.players[1 - aMatch.playerID].timestampSignature, aMatch.opponentSubkeySignature);

    if (moveMaker(metaState, cheaterMove, aMatch.opponentSubkeySignature) != opponent) {
      return false;
    }

    (isLegal,) = aMatch.game.isMoveLegal(metaState, cheaterMove.move);

    if (isLegal) {
      return false;
    }

    return true;
  }

  // XXX: https://github.com/ethereum/solidity/issues/3199#issuecomment-365035663
  function reportCheater(Match aMatch, DGame.MetaState metaState, Move cheaterMove) public {
    address opponent;
    uint value;

    require(canReportCheater(aMatch, metaState, cheaterMove));

    invalidateTimestamp(aMatch.timestamp, aMatch.players[0].timestampSignature);
    invalidateTimestamp(aMatch.timestamp, aMatch.players[1].timestampSignature);

    opponent = playerAccount(aMatch.timestamp, aMatch.players[1 - aMatch.playerID].timestampSignature, aMatch.opponentSubkeySignature);
    value = balance[opponent];
    delete balance[opponent];
    balance[msg.sender] += value / 2;
    balance[owner] += (value + 1) / 2;

    aMatch.game.registerCheat(opponent);

    emit balanceChanged(opponent);
    emit balanceChanged(msg.sender);
    emit balanceChanged(owner);
    emit cheaterReported(opponent, timestampSubkey(aMatch.timestamp, aMatch.players[1 - aMatch.playerID].timestampSignature), aMatch.timestamp);
  }

  event rewardClaimed(address indexed account, address indexed subkey, uint indexed timestamp);
  event cheaterReported(address indexed account, address indexed subkey, uint indexed timestamp);

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

  // XXX: abigen: Failed to generate ABI binding: unsupported arg type: tuple
  function subkeyParentXXX(address subkey, uint8 subkeyV, bytes32 subkeyR, bytes32 subkeyS) public pure returns (address) {
    return subkeyParent(subkey, SubkeySignature(subkeyV, subkeyR, subkeyS));
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

  function timestampSubkeyXXX(uint timestamp, uint8 timestampV, bytes32 timestampR, bytes32 timestampS) public pure returns (address) {
    return timestampSubkey(timestamp, TimestampSignature(timestampV, timestampR, timestampS));
  }

  function timestampSubkey(uint timestamp, TimestampSignature timestampSignature) public pure returns (address) {
    return ecrecover(keccak256(timestamp), timestampSignature.v, timestampSignature.r, timestampSignature.s);
  }

  // XXX: abigen: Failed to generate ABI binding: unsupported arg type: tuple
  function playerAccountXXX(uint timestamp, uint8 timestampV, bytes32 timestampR, bytes32 timestampS, uint8 subkeyV, bytes32 subkeyR, bytes32 subkeyS) public pure returns (address) {
    return playerAccount(timestamp, TimestampSignature(timestampV, timestampR, timestampS), SubkeySignature(subkeyV, subkeyR, subkeyS));
  }

  // XXX: https://github.com/ethereum/solidity/issues/3275#issuecomment-365087323
  function playerAccount(uint timestamp, TimestampSignature timestampSignature, SubkeySignature subkeySignature) public pure returns (address) {
    return subkeyParent(timestampSubkey(timestamp, timestampSignature), subkeySignature);
  }

  // XXX: https://github.com/ethereum/solidity/issues/3270
  function matchHash(DGame game, uint timestamp, address[2] accounts, address[2] subkeys, uint32[2] seedRatings, bytes32[PUBLIC_SEED_LENGTH][2] publicSeeds) public pure returns (bytes32) {
    return keccak256(game, timestamp, accounts[0], accounts[1], subkeys[0], subkeys[1], seedRatings[0], seedRatings[1], publicSeeds[0], publicSeeds[1]);
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

  function invalidateTimestamp(uint timestamp, TimestampSignature timestampSignature) private {
    invalidatedTimestamps[keccak256(timestamp, timestampSubkey(timestamp, timestampSignature))] = true;
  }

  function isTimestampInvalid(uint timestamp, TimestampSignature timestampSignature) private view returns (bool) {
    return invalidatedTimestamps[keccak256(timestamp, timestampSubkey(timestamp, timestampSignature))];
  }

  function matchMaker(Match aMatch, address sender) private pure returns (address) {
    address[2] memory accounts;
    address[2] memory subkeys;
    uint32[2] memory seedRatings;
    // XXX: https://github.com/ethereum/solidity/issues/3270
    bytes32[PUBLIC_SEED_LENGTH][2] memory publicSeeds;
    bytes32 hash;

    accounts[aMatch.playerID] = sender;
    accounts[1 - aMatch.playerID] = playerAccount(aMatch.timestamp, aMatch.players[1 - aMatch.playerID].timestampSignature, aMatch.opponentSubkeySignature);
    subkeys[0] = timestampSubkey(aMatch.timestamp, aMatch.players[0].timestampSignature);
    subkeys[1] = timestampSubkey(aMatch.timestamp, aMatch.players[1].timestampSignature);
    seedRatings[0] = aMatch.players[0].seedRating;
    seedRatings[1] = aMatch.players[1].seedRating;
    publicSeeds[0] = aMatch.players[0].publicSeed;
    publicSeeds[1] = aMatch.players[1].publicSeed;
    hash = matchHash(aMatch.game, aMatch.timestamp, accounts, subkeys, seedRatings, publicSeeds);

    return ecrecover(hash, aMatch.matchSignature.v, aMatch.matchSignature.r, aMatch.matchSignature.s);
  }

  address private owner;
  mapping(bytes32 => bool) private invalidatedTimestamps;
}
