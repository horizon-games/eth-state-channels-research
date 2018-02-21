pragma solidity ^0.4.19;
pragma experimental ABIEncoderV2;

import './DGame.sol';

contract Arcadeum {
  uint constant openChannelRestrictedGas = 21000; // XXX
  uint constant closeChannelRestrictedGas = 21000; // XXX
  uint constant playerWonRestrictedGas = 21000; // XXX
  uint constant playerCheatedRestrictedGas = 21000; // XXX

  function Arcadeum() public {
    owner = msg.sender;
  }

  function deposit() public payable {
    balance[msg.sender] += msg.value;
  }

  function withdraw(uint amount) public {
    require(amount <= balance[msg.sender]);

    balance[msg.sender] -= amount;
    msg.sender.transfer(amount);
  }

  function didPlayerWin(DGame.Match dMatch, DGame.MetaState mState, DGame.MetaMove loserMove, DGame.MetaMove[] winnerMoves) public view returns (bool) {
    bytes24 channelID;
    uint loserID;
    uint winnerID;
    uint winner;
    uint nextPlayers;
    DGame.MetaState memory nextState;
    uint i;

    channelID = (bytes24(address(dMatch.game)) << 32) | bytes24(dMatch.matchID);

    if (channel[channelID].stake == 0) {
      return false;
    }

    if (matchMaker(dMatch) != owner) {
      return false;
    }

    loserID = loserMove.move.playerID;
    winnerID = 1 - loserID;

    if (!dMatch.game.isSubkeySigned(dMatch, loserID)) {
      return false;
    }

    if (moveMaker(dMatch, mState, loserMove) != dMatch.players[loserID].subkey) {
      return false;
    }

    winner = dMatch.game.winner(mState);

    if (winner == winnerID + 1) {
      return true;
    }

    nextPlayers = dMatch.game.nextPlayers(mState);

    /* XXX: https://github.com/ethereum/solidity/issues/3516 *

    if (nextPlayers == 3) {
      if (winnerMoves[0].move.playerID != winnerID) {
        return false;
      }

      nextState = dMatch.game.nextState2(mState, loserMove.move, winnerMoves[0].move);

      i = 1;

    } else {
      if (nextPlayers != loserID + 1) {
        return false;
      }

      nextState = dMatch.game.nextState1(mState, loserMove.move);

      i = 0;

    }

    for (; i < winnerMoves.length; i++) {
      if (winnerMoves[i].move.playerID != winnerID) {
        return false;
      }

      nextPlayers = dMatch.game.nextPlayers(nextState);

      if (nextPlayers != winnerID + 1) {
        return false;
      }

      nextState = dMatch.game.nextState1(nextState, winnerMoves[i].move);
    }

    * XXX: https://github.com/ethereum/solidity/issues/3516 */

    winner = dMatch.game.winner(nextState);

    return winner == winnerID + 1;
  }

  function didPlayerCheat(DGame.Match dMatch, DGame.MetaState mState, DGame.MetaMove cheaterMove) public view returns (bool) {
    bytes24 channelID;
    uint cheaterID;

    channelID = (bytes24(address(dMatch.game)) << 32) | bytes24(dMatch.matchID);

    if (channel[channelID].stake == 0) {
      return false;
    }

    if (matchMaker(dMatch) != owner) {
      return false;
    }

    cheaterID = cheaterMove.move.playerID;

    if (!dMatch.game.isSubkeySigned(dMatch, cheaterID)) {
      return false;
    }

    if (moveMaker(dMatch, mState, cheaterMove) != dMatch.players[cheaterID].subkey) {
      return false;
    }

    if (dMatch.game.isMoveLegal(mState, cheaterMove.move)) {
      return false;
    }

    return true;
  }

  function playerWon(DGame.Match dMatch, DGame.MetaState mState, DGame.MetaMove loserMove, DGame.MetaMove[] winnerMoves) public {
    require(didPlayerWin(dMatch, mState, loserMove, winnerMoves));

    playerWonInternal(dMatch.game, dMatch.matchID, 1 - loserMove.move.playerID, 0);
  }

  function playerCheated(DGame.Match dMatch, DGame.MetaState mState, DGame.MetaMove cheaterMove) public {
    require(didPlayerCheat(dMatch, mState, cheaterMove));

    playerCheatedInternal(dMatch.game, dMatch.matchID, cheaterMove.move.playerID, 0);
  }

  modifier restricted() { require(msg.sender == owner); _; }

  function openChannelRestricted(DGame game, uint32 matchID, uint stake, address[2] players) public restricted {
    bytes24 channelID;

    channelID = (bytes24(address(game)) << 32) | bytes24(matchID);

    require(stake > 0);
    require(balance[players[0]] >= stake);
    require(balance[players[1]] >= stake);
    require(channel[channelID].stake == 0);

    balance[players[0]] -= stake;
    balance[players[1]] -= stake;

    channel[channelID].stake = stake;
    channel[channelID].openCost = (openChannelRestrictedGas * tx.gasprice + 1) & ~uint(1);
    channel[channelID].players = players;
  }

  function closeChannelRestricted(DGame game, uint32 matchID) public restricted {
    bytes24 channelID;
    uint closeCost;
    uint refund;

    channelID = (bytes24(address(game)) << 32) | bytes24(matchID);

    require(channel[channelID].stake > 0);

    closeCost = (closeChannelRestrictedGas * tx.gasprice + 1) & ~uint(1);
    refund = channel[channelID].stake - channel[channelID].openCost / 2 - closeCost / 2;

    balance[channel[channelID].players[0]] += refund;
    balance[channel[channelID].players[1]] += refund;
    balance[owner] += channel[channelID].openCost + closeCost;
    delete channel[channelID];
  }

  function playerWonRestricted(DGame game, uint32 matchID, uint winnerID) public restricted {
    bytes24 channelID;

    channelID = (bytes24(address(game)) << 32) | bytes24(matchID);

    require(channel[channelID].stake > 0);

    playerWonInternal(game, matchID, winnerID, playerWonRestrictedGas);
  }

  function playerCheatedRestricted(DGame game, uint32 matchID, uint cheaterID) public restricted {
    bytes24 channelID;

    channelID = (bytes24(address(game)) << 32) | bytes24(matchID);

    require(channel[channelID].stake > 0);

    playerCheatedInternal(game, matchID, cheaterID, playerCheatedRestrictedGas);
  }

  function playerWonInternal(DGame game, uint32 matchID, uint winnerID, uint closeGas) private {
    bytes24 channelID;
    uint closeCost;

    channelID = (bytes24(address(game)) << 32) | bytes24(matchID);
    closeCost = closeGas * tx.gasprice;

    balance[channel[channelID].players[winnerID]] += channel[channelID].stake - channel[channelID].openCost - closeCost;
    balance[channel[channelID].players[1 - winnerID]] += channel[channelID].stake;
    balance[owner] += channel[channelID].openCost + closeCost;
    delete channel[channelID];
  }

  function playerCheatedInternal(DGame game, uint32 matchID, uint cheaterID, uint closeGas) private {
    bytes24 channelID;
    uint closeCost;

    channelID = (bytes24(address(game)) << 32) | bytes24(matchID);
    closeCost = closeGas * tx.gasprice;

    balance[channel[channelID].players[1 - cheaterID]] += channel[channelID].stake * 3 / 2 - channel[channelID].openCost - closeCost;
    balance[owner] += (channel[channelID].stake + 1) / 2 + channel[channelID].openCost + closeCost;
    delete channel[channelID];
  }

  function matchMaker(DGame.Match dMatch) private pure returns (address) {
    bytes32 hash;

    hash = keccak256(dMatch.game, dMatch.matchID, dMatch.players[0].account, dMatch.players[1].account, dMatch.players[0].publicSeed, dMatch.players[1].publicSeed);

    return ecrecover(hash, dMatch.signature.v, dMatch.signature.r, dMatch.signature.s);
  }

  function moveMaker(DGame.Match dMatch, DGame.MetaState mState, DGame.MetaMove mMove) private pure returns (address) {
    bytes32 matchHash;
    bytes32 stateHash;
    bytes32 moveHash;
    bytes32 hash;

    matchHash = keccak256(dMatch.game, dMatch.matchID, dMatch.players[0].account, dMatch.players[1].account);
    stateHash = keccak256(mState.nonce, mState.tag, mState.data, mState.state.tag, mState.state.data);
    moveHash = keccak256(mMove.move.playerID, mMove.move.data);
    hash = keccak256(matchHash, stateHash, moveHash);

    return ecrecover(hash, mMove.signature.v, mMove.signature.r, mMove.signature.s);
  }

  address owner;
  mapping(address => uint) balance;
  mapping(bytes24 => Channel) channel;

  struct Channel {
    uint stake;
    uint openCost;
    address[2] players;
  }
}
