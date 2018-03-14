package matcher

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/asn1"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/websocket"
	"github.com/horizon-games/arcadeum/server/config"
	"github.com/horizon-games/arcadeum/server/services/arcadeum"
	cr "github.com/horizon-games/arcadeum/server/services/crypto"
	"github.com/horizon-games/arcadeum/server/services/util"
)

type Code int
type Status int

const (
	ERROR            Code = -1
	MSG              Code = 0 // code for messages passed between players during gameplay
	INIT             Code = 1 // match found
	SIGNED_TIMESTAMP Code = 2
	MATCH_VERIFIED   Code = 3 // all players in match session have passed all validation tests
)

const (
	Unknown      Status = 0
	Waiting      Status = 1
	Moving       Status = 2
	Won          Status = 3
	Lost         Status = 4
	Disqualified Status = 5
)

type MatchRequest struct {
	Conn *websocket.Conn
	*Token
}

type MatchResponse struct {
	Account common.Address // Owner of seed deck; this value is derived
	Rank    uint32         // calculated rank of player based on seed "deck"
	Request *MatchRequest
}

type PlayerInfo struct {
	Rank  uint32 // Player rank as returned from Solidity pure function
	Index uint8  // index in match session; arbitrarily set when match found and player joins a session
	Conn  *websocket.Conn
	*Token

	SeedHash     []byte         // Hash of seed as returned from Solidity pure function
	Account      common.Address // owner account of signed subkey (derived from Token); See Token
	TimestampSig cr.Signature   // TimestampSig of session timestamp by this player
	Verified     bool           // true if the player has proven their TimestampSig
}

// Represents a matched game session
type Session struct {
	GameID    uint32
	MatchID   uint32
	Player1   *PlayerInfo
	Player2   *PlayerInfo
	Timestamp int64         // Game start in Unix time
	Signature *cr.Signature // Matcher's session signature
}

type InitMessage struct {
	MatchID   uint32 `json:"matchID"`
	Timestamp int64  `json:"timestamp"`
}

type Meta struct {
	MatchID uint32 `json:"matchID"` // uuid
	Index   uint8  `json:"index"`   // index of player in game, i.e., player ID
	Code    Code   `json:"code"`    // message type
}

// Message send over the wire between players
type Message struct {
	Meta    `json:"meta"`
	Payload string `json:"payload"`
}

type Service struct {
	MatchID   uint32 // auto-incrementing match count
	ArcClient *arcadeum.Client

	// gameID -> map rank -> account -> player
	WaitPool    map[uint32]map[uint32]map[common.Address]*MatchResponse
	SessionPool map[uint32]*Session // map matchID -> in-game session; matchIDs are globally unique

	ENV    *config.ENVConfig
	Config *config.MatcherConfig
}

var matchResponseChannel = make(chan *MatchResponse)

func NewService(
	env *config.ENVConfig,
	cfg *config.MatcherConfig,
	ethcfg *config.ETHConfig,
	arcconfig *config.ArcadeumConfig) *Service {
	service := &Service{
		MatchID:     0,
		WaitPool:    make(map[uint32]map[uint32]map[common.Address]*MatchResponse),
		SessionPool: make(map[uint32]*Session),
		ENV:         env,
		Config:      cfg,
		ArcClient:   arcadeum.NewArcadeumClient(ethcfg, arcconfig),
	}
	go service.ArcClient.HandleWithdrawalStarted(service)
	return service
}

// Event handler when we have detected when a user has decided to withdraw money from their account
func (s *Service) OnWithdrawalStarted(event *arcadeum.ArcadeumWithdrawalStarted) {
	account := event.Account
	sess := s.FindSessionByAccount(account)
	gameaddr := s.ArcClient.GameAddress[sess.GameID]
	contract := s.ArcClient.ArcadeumContract

	withdrawing, err := contract.IsWithdrawing(&bind.CallOpts{}, account)
	if err != nil {
		log.Println("ERROR: could not verify IsWithdrawing state", err)
		return
	}
	if withdrawing {
		return //! how to get notified when withdrawal complete?
	}

	player, err := sess.FindPlayerByAccount(account)
	if err != nil {
		log.Printf("ERROR: could not find account %s in session", account)
		return
	}

	var playerR, playerS, sessR, sessS [32]byte
	copy(playerR[:], player.TimestampSig.R)
	copy(playerS[:], player.TimestampSig.S)
	copy(sessR[:], sess.Signature.R)
	copy(sessS[:], sess.Signature.S)
	canWithdraw, err := contract.CanStopWithdrawalXXX(
		&bind.CallOpts{},
		gameaddr,
		sess.MatchID,
		big.NewInt(sess.Timestamp),
		player.TimestampSig.V,
		playerR,
		playerS,
		sess.Signature.V,
		sessR,
		sessS)
	if err != nil {
		log.Printf("ERROR: Could not read CanStopWithdrawal() value from blockchain", err)
		return
	}
	if !canWithdraw { // Slash player
		opts := s.NewKeyedTransactor()
		opts.From = s.Config.AccountAddress
		opts.Value = nil    // no funds
		opts.GasLimit = 0   // estimate
		opts.GasPrice = nil // use price oracle
		_, err := contract.StopWithdrawalXXX(
			opts,
			gameaddr,
			sess.MatchID,
			big.NewInt(sess.Timestamp),
			player.TimestampSig.V,
			playerR,
			playerS,
			sess.Signature.V,
			sessR,
			sessS)
		if err != nil {
			log.Printf("ERROR: failure to slash withdrawal matchID %d account %s", sess.MatchID, player.Account)
			return
		}
	}
}

func (s *Session) Rank() uint32 {
	return s.Player1.Rank // both players in session have same rank so just return first one
}

func (s *Session) FindPlayer(ws *websocket.Conn) *PlayerInfo {
	if s.Player1.Conn == ws {
		return s.Player1
	}
	return s.Player2
}

func (s *Session) FindPlayerByAccount(account common.Address) (*PlayerInfo, error) {
	if s.Player1.Account == account {
		return s.Player1, nil
	} else if s.Player2.Account == account {
		return s.Player2, nil
	}
	return nil, errors.New("Unknonw account " + account.String())
}

func (s *Session) IsVerified() bool {
	return s.Player1.Verified && s.Player2.Verified
}

func (s *Session) GetOpponent(ws *websocket.Conn) *PlayerInfo {
	if s.Player1.Conn == ws {
		return s.Player2
	}
	return s.Player1
}

func (s *Session) ContainsAccount(addr common.Address) bool {
	return s.Player1.Account == addr || s.Player2.Account == addr
}

func (s *Service) VerifyTimestamp(gameID uint32, matchID uint32, req *arcadeum.VerifyTimestampRequest, player *PlayerInfo) (bool, error) {
	account, err := s.ArcClient.VerifySignedTimestamp(gameID, matchID, req, player.SubKeySignature)
	if err != nil {
		return false, errors.New("Could not deserialize signed timestamp payload.")
	}
	return player.Account == account, nil
}

func Context(r *http.Request) *Token {
	return r.Context().Value("Token").(*Token)
}

func NewError(message string) Message {
	return Message{Meta: Meta{Code: ERROR}, Payload: message}
}

func (s *Service) FindMatch(req *MatchRequest) {
	response, err := s.Authenticate(req)
	if err != nil {
		log.Println("Error authenticating match request. Closing connection.", err)
		req.Conn.WriteJSON(NewError(err.Error()))
		req.Conn.Close()
		return
	}
	matchResponseChannel <- response
}

func (s *Service) Authenticate(req *MatchRequest) (*MatchResponse, error) {
	address, err := s.ArcClient.SubKeyParent(req.SubKey, req.SubKeySignature)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error validating subkey account address. %s", err.Error()))
	}
	status, err := s.ArcClient.GetStakedStatus(address)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error validating stake.", err))
	}
	if status == arcadeum.STAKED {
		owner, err := s.ArcClient.IsSecretSeedValid(req.GameID, address, req.Seed)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error verifying seed ownership.", err))
		}
		if !owner {
			return nil, errors.New("Invalid seed ownership.")
		}
		rank, err := s.ArcClient.CalculateRank(req.GameID, req.Seed)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error calculating rank.", err))
		}
		return &MatchResponse{
			Account: address,
			Rank:    rank,
			Request: req,
		}, nil
	} else if status == arcadeum.STAKED_INSUFFICIENT_BALANCE {
		return nil, errors.New("Insufficient stake balance.")
	} else {
		return nil, errors.New("Player has not staked.")
	}
}

func (s *Service) HandleMatchResponses() {
	for {
		rp := <-matchResponseChannel
		s.Match(rp)
	}
}

func (s *Service) Match(rp *MatchResponse) {
	opponent := s.FindMatchByRank(rp.Request.GameID, rp.Rank)
	if opponent != nil {
		log.Printf("Match found!")
		s.InitializeGame(rp, opponent)
	} else {
		log.Println("Match not found, adding to match pool")
		s.AddToMatchPool(rp)
	}
}

func (s *Service) RemoveFromWaitingPool(resps ...*MatchResponse) {
	for _, r := range resps {
		delete(s.WaitPool[r.Request.GameID][r.Rank], r.Account)
	}
}

// We didn't find a match of equal rank, so add to wait pool
func (s *Service) AddToMatchPool(rp *MatchResponse) {
	gameID := rp.Request.GameID
	rank := rp.Rank
	if s.WaitPool[gameID] == nil {
		s.WaitPool[gameID] = make(map[uint32]map[common.Address]*MatchResponse)
	}
	if s.WaitPool[gameID][rank] == nil {
		s.WaitPool[gameID][rank] = make(map[common.Address]*MatchResponse)
	}
	s.WaitPool[gameID][rank][rp.Account] = rp
}

func (s *Service) InitializeGame(p1 *MatchResponse, p2 *MatchResponse) {
	session, err := s.CreateSession(p1, p2)
	if err != nil {
		log.Println("error creating session: ", err)
		s.Close("Error creating match session. Closing match connection.", p1, p2)
		return
	}
	s.RequestTimestampProof(session)
	s.AddToSessionPool(session)
}

func (s *Service) Close(message string, p1 *MatchResponse, p2 *MatchResponse) {
	p1.Request.Conn.WriteJSON(NewError(message))
	p2.Request.Conn.WriteJSON(NewError(message))
	s.RemoveFromWaitingPool(p1, p2)
	p1.Request.Conn.Close()
	p2.Request.Conn.Close()
}

// Session has been verified so begin match
func (srv *Service) BeginVerifiedMatch(s *Session) error {
	if !s.IsVerified() {
		return nil
	}
	msg, err := srv.BuildMatchVerifiedMessageWithSignature(s)
	if err != nil {
		return err
	}
	msg.SignatureOpponentTimestamp = &s.Player2.TimestampSig
	msg.SignatureOpponentSubkey = &s.Player2.Token.SubKeySignature
	s.Signature = msg.SignatureMatchHash
	relaymsg := &Message{
		Meta: Meta{
			MatchID: s.MatchID,
			Code:    MATCH_VERIFIED,
		},
		Payload: util.Jsonify(msg),
	}
	err2 := s.Player1.Conn.WriteJSON(relaymsg)
	if err2 != nil {
		return err2
	}
	msg.PlayerIndex = 1
	msg.SignatureOpponentTimestamp = &s.Player1.TimestampSig
	msg.SignatureOpponentSubkey = &s.Player1.Token.SubKeySignature
	relaymsg = &Message{
		Meta: Meta{
			MatchID: s.MatchID,
			Code:    MATCH_VERIFIED,
		},
		Payload: util.Jsonify(msg),
	}
	err3 := s.Player2.Conn.WriteJSON(relaymsg)
	if err3 != nil {
		return err3
	}
	return nil
}

func (s *Service) NewKeyedTransactor() *bind.TransactOpts {
	privkey := s.PrivKey()
	return bind.NewKeyedTransactor(privkey)
}

func (srv *Service) BuildMatchVerifiedMessageWithSignature(s *Session) (*arcadeum.MatchVerifiedMessage, error) {
	msg := &arcadeum.MatchVerifiedMessage{
		Accounts:    [2]common.Address{s.Player1.Account, s.Player2.Account},
		GameAddress: srv.ArcClient.GameAddress[s.GameID],
		MatchID:     s.MatchID,
		Timestamp:   s.Timestamp,
		Players: [2]*arcadeum.MatchVerifiedPlayerInfo{
			{
				SeedRating: s.Player1.Rank,
				PublicSeed: s.Player1.SeedHash,
			},
			{
				SeedRating: s.Player2.Rank,
				PublicSeed: s.Player2.SeedHash,
			},
		},
	}
	hash, err := srv.ArcClient.MatchHash(msg)
	if err != nil {
		return nil, err
	}
	msg.MatchHash = hash

	// Have the matcher sign
	sig, err := crypto.Sign(hash[:], srv.PrivKey())
	if err != nil {
		return nil, err
	}
	msg.SignatureMatchHash = &cr.Signature{
		V: 27 + sig[64],
		R: sig[0:32],
		S: sig[32:64],
	}

	return msg, nil
}

func (srv *Service) RequestTimestampProof(s *Session) error {
	log.Println("Requesting timestamp proof from both players")
	err := WriteInitMessage(s.MatchID, s.Timestamp, s.Player1)
	if err != nil {
		return err
	}
	err2 := WriteInitMessage(s.MatchID, s.Timestamp, s.Player2)
	if err2 != nil {
		return err2
	}
	return nil
}

func WriteInitMessage(matchID uint32, timestamp int64, p *PlayerInfo) error {
	payload := &InitMessage{
		MatchID:   matchID,
		Timestamp: timestamp,
	}
	return p.Conn.WriteJSON(
		Message{
			Meta:    Meta{MatchID: matchID, Code: INIT, Index: p.Index},
			Payload: util.Jsonify(payload)})
}

func (s *Service) CreateSession(p1 *MatchResponse, p2 *MatchResponse) (*Session, error) {
	s.MatchID += 1 //! make atomic
	player1, err := s.BuildPlayerInfo(p1)
	if err != nil {
		return nil, err
	}
	player2, err := s.BuildPlayerInfo(p2)
	if err != nil {
		return nil, err
	}
	RandomizeTurn(player1, player2)
	if err != nil {
		return nil, err
	}
	gameID := p1.Request.GameID // arbitrarily choose first player to get game ID
	duration, err := s.ArcClient.MatchDuration(gameID)
	if err != nil {
		return nil, err
	}
	return &Session{
		GameID:    gameID,
		MatchID:   s.MatchID,
		Player1:   player1,
		Player2:   player2,
		Timestamp: time.Now().Add(duration).Unix(),
	}, nil
}

func RandomizeTurn(p1 *PlayerInfo, p2 *PlayerInfo) {
	p1.Index = 0
	p2.Index = 1
}

func (srv *Service) SignElliptic(inputs ...interface{}) (r, s *big.Int, err error) {
	compact, err := Compact(inputs...)
	if err != nil {
		return
	}
	hash := crypto.Keccak256(compact)
	privkey := srv.PrivKey()
	r, s, err = ecdsa.Sign(rand.Reader, privkey, hash)
	return
}

func (srv *Service) PrivKey() *ecdsa.PrivateKey {
	path := fmt.Sprintf("%s/%s", srv.ENV.WorkingDir, srv.Config.PrivKeyFile)
	privkey, err := crypto.LoadECDSA(path)
	if err != nil {
		log.Fatalf("Invalid private key")
	}
	return privkey
}

func (srv *Service) Sign(inputs ...interface{}) ([]byte, error) {
	r, s, err := srv.SignElliptic(inputs...)
	if err != nil {
		return nil, err
	}
	return asn1.Marshal(cr.EcdsaSignature{r, s})
}

func Compact(inputs ...interface{}) ([]byte, error) {
	var compact []byte
	for _, elem := range inputs {
		b, err := IToB(elem)
		if err != nil {
			return nil, err
		}
		compact = append(compact, b...)
	}
	return compact, nil
}

func IToB(data interface{}) ([]byte, error) {
	if _, ok := data.(string); ok {
		return []byte(data.(string)), nil
	}
	if _, ok := data.([]byte); ok {
		return data.([]byte), nil
	}
	if s, ok := data.(int); ok {
		data = uint32(s)
	}
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *Service) BuildPlayerInfo(p *MatchResponse) (*PlayerInfo, error) {
	seedHash, err := s.ArcClient.PublicSeed(p.Request.GameID, p.Request.Seed)
	if err != nil {
		return nil, err
	}
	return &PlayerInfo{
		Rank:     p.Rank,
		Conn:     p.Request.Conn,
		Token:    p.Request.Token,
		SeedHash: seedHash,
		Account:  p.Account,
	}, nil
}

func (s *Service) FindSession(matchID uint32) *Session {
	return s.SessionPool[matchID]
}

func (s *Service) FindSessionByAccount(account common.Address) *Session {
	for _, s := range s.SessionPool {
		if s.ContainsAccount(account) {
			return s
		}
	}
	return nil
}

// Find a similarly ranked player to play against.
// This should be done in a random, non-forgeable way.
func (s *Service) FindMatchByRank(gameID uint32, rank uint32) *MatchResponse { //! make atomic
	waiting := s.WaitPool[gameID][rank]
	if len(waiting) > 0 {
		var resp *MatchResponse
		// Choose the first element in the map (we an make this better)
		for _, r := range waiting {
			// one solution: randomly select index and return when you git that index
			resp = r
			break
		}
		delete(waiting, resp.Account)
		return resp
	}

	return nil
}

func (s *Service) AddToSessionPool(session *Session) { //! make atomic
	waiting := s.WaitPool[session.GameID][session.Rank()]
	delete(waiting, session.Player1.Account)
	delete(waiting, session.Player2.Account)
	s.SessionPool[session.MatchID] = session
}
