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

	"encoding/json"
	"strconv"

	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/horizon-games/arcadeum/server/config"
	"github.com/horizon-games/arcadeum/server/lib/arcadeum"
	"github.com/horizon-games/arcadeum/server/lib/crypto"
	"github.com/horizon-games/arcadeum/server/lib/util"
	"github.com/pborman/uuid"
)

type Code int
type Status int

const (
	TERMINATE Code = -2 // unrecoverable fatal error
	ERROR     Code = -1

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

type MatchResponse struct {
	Account common.Address // Owner of seed deck; this value is derived
	Rank    uint32         // calculated rank of player based on seed "deck"
	*Token
}

type InitMessage struct {
	Timestamp int64 `json:"timestamp"`
}

type Meta struct {
	Index  uint8           `json:"index"` // index of player in game, i.e., player ID
	Code   Code            `json:"code"`  // message type
	SubKey *common.Address `json:"subkey"`
}

// Message send over the wire between players
type Message struct {
	*Meta   `json:"meta"`
	Payload string `json:"payload"`
}

type Service struct {
	ArcClient *arcadeum.Client
	ENV       *config.ENVConfig
	Config    *config.MatcherConfig
	*SessionManager
	*PubSubManager
}

// Responses in this channel have already been authenticated
var matchResponseChannel = make(chan *MatchResponse)

func NewService(
	env *config.ENVConfig,
	cfg *config.MatcherConfig,
	ethcfg *config.ETHConfig,
	arcconfig *config.ArcadeumConfig,
	rediscfg *config.RedisConfig) *Service {
	sessMgr := NewSessionManager(rediscfg)
	service := &Service{
		ENV:            env,
		Config:         cfg,
		ArcClient:      arcadeum.NewArcadeumClient(ethcfg, arcconfig),
		SessionManager: sessMgr,
		PubSubManager:  NewPubSubManager(sessMgr),
	}
	go service.ArcClient.HandleWithdrawalStarted(service)
	return service
}

// Event handler when we have detected when a user has decided to withdraw money from their account
func (s *Service) OnWithdrawalStarted(event *arcadeum.ArcadeumWithdrawalStarted) {
	account := event.Account
	sess, err := s.GetSessionByAccount(&account)
	if err != nil {
		log.Println("ERROR: Could not find session %s", err.Error())
		return
	}
	if sess.IsEmpty() {
		log.Println("ERROR: empty session for account %s", account.String())
		return
	}
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

	canWithdraw, err := contract.CanStopWithdrawalXXX(
		&bind.CallOpts{},
		big.NewInt(sess.Timestamp),
		player.TimestampSig.V,
		player.TimestampSig.R,
		player.TimestampSig.S,
		sess.Signature.V,
		sess.Signature.R,
		sess.Signature.S)
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
			big.NewInt(sess.Timestamp),
			player.TimestampSig.V,
			player.TimestampSig.R,
			player.TimestampSig.S,
			sess.Signature.V,
			sess.Signature.R,
			sess.Signature.S)
		if err != nil {
			log.Printf("ERROR: failure to slash withdrawal account %s", player.Account)
			return
		}
	}
}

func (s *Session) Rank() uint32 {
	return s.Player1.Rank // both players in session have same rank so just return first one
}

func (s *Session) IsVerified() bool {
	return s.Player1.Verified && s.Player2.Verified
}

func (s *Session) FindPlayerBySubKey(subKey *common.Address) *PlayerInfo {
	if s.Player1.SubKey.String() == subKey.String() {
		return s.Player1
	} else if s.Player2.SubKey.String() == subKey.String() {
		return s.Player2
	} else {
		return nil
	}
}

func (s *Session) IsWaiting() bool {
	return s.Player1 == nil || s.Player2 == nil
}

func (s *Session) FindOpponent(subKey *common.Address) *PlayerInfo {
	if s.Player1 != nil && s.Player1.SubKey.String() == subKey.String() {
		return s.Player2
	} else if s.Player2 != nil && s.Player2.SubKey.String() == subKey.String() {
		return s.Player1
	}
	return nil
}

func (s *Service) VerifyTimestamp(req *arcadeum.VerifyTimestampRequest, player *PlayerInfo) (bool, error) {
	account, err := s.ArcClient.VerifySignedTimestamp(req, player.SubKeySignature)
	if err != nil {
		return false, errors.New("Could not deserialize signed timestamp payload.")
	}
	return *player.Account == account, nil
}

func Context(r *http.Request) *Token {
	return r.Context().Value("Token").(*Token)
}

func NewError(message string) Message {
	return Message{Meta: &Meta{Code: ERROR}, Payload: message}
}

func NewTerminateMessage(message string) Message {
	return Message{Meta: &Meta{Code: TERMINATE}, Payload: message}
}

func (s *Service) Reconnect(token *Token) (bool, error) {
	sess, err := s.GetSessionBySubKey(token.SubKey)
	if err != nil {
		return false, err
	}
	if sess.IsEmpty() {
		return false, nil
	}
	if sess.IsWaiting() {
		return true, nil
	} else { // session has already matched
		p := sess.FindPlayerBySubKey(token.SubKey)
		if !p.Verified {
			err := s.SendTimestampProof(p, sess.Timestamp) // ask again
			return true, err
		}
	}
	return true, nil
}

func (s *Service) FindMatch(token *Token) {
	reconnection, err := s.Reconnect(token)
	if err != nil {
		s.PublishToSubKey(token.SubKey, NewTerminateMessage(fmt.Sprintf("failure finding match, disconnecting", err.Error())))
		return
	}
	if reconnection {
		return // behave as if you never disconnected
	}
	response, err := s.Authenticate(token)
	if err != nil {
		message := fmt.Sprintf("Error authenticating match request. Closing connection. %s", err.Error())
		s.PublishToSubKey(token.SubKey, NewTerminateMessage(message))
		return
	}
	matchResponseChannel <- response
}

func (s *Service) Authenticate(token *Token) (*MatchResponse, error) {
	address, err := s.ArcClient.SubKeyParent(*token.SubKey, token.SubKeySignature)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error validating subkey account address. %s", err.Error()))
	}
	status, err := s.ArcClient.GetStakedStatus(address)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error validating stake.", err))
	}
	if status == arcadeum.STAKED {
		log.Println("Seed information: GameID: %d, address: %s, seed byte length: %d, seed: %s", token.GameID, address, len(token.Seed), string(token.Seed))
		owner, err := s.ArcClient.IsSecretSeedValid(token.GameID, address, token.Seed)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error verifying seed ownership.", err))
		}
		if !owner {
			return nil, errors.New("Invalid seed ownership.")
		}
		rank, err := s.ArcClient.CalculateRank(token.GameID, token.Seed)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error calculating rank.", err))
		}
		return &MatchResponse{
			Account: address,
			Rank:    rank,
			Token:   token,
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

// Invariant: rp has been authenticated
func (s *Service) Match(rp *MatchResponse) {
	uuid, err := s.TakeRandomSessionByRank(rp.Rank)
	if err != nil {
		s.Close(fmt.Sprintf("Error finding opponent %s", err.Error()), rp)
		return
	}
	if uuid.IsEmpty() {
		err = s.AddToMatchPool(rp)
	} else {
		err = s.InitGame(uuid, rp)
	}
	if err != nil {
		s.ReaddToMatchPool(rp.Rank, uuid)
		s.Close(fmt.Sprintf("Error looking for match %s", err.Error()), rp)
	}
}

func (s *Service) InitGame(uid UUID, r *MatchResponse) error {
	session, err := s.GetSessionByID(uid)
	if err != nil {
		return err
	}
	if session.IsEmpty() {
		return errors.New("Trying to match with empty session")
	}
	newSess, err := s.CreateSession(r)
	if err != nil {
		return err
	}
	newSess.Player1.Index = 1
	session.Player2 = newSess.Player1
	session.Timestamp = time.Now().Unix()
	err = s.UpdateSession(session)
	if err != nil {
		return err
	}
	return s.RequestTimestampProof(session)
}

func (s *Service) AddToMatchPool(r *MatchResponse) error {
	session, err := s.CreateSession(r)
	if err != nil {
		return err
	}
	err = s.SessionManager.AddToMatchPool(session)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) RemoveFromWaitingPool(resps ...*MatchResponse) {
	//! Cleanup session cache
}

func (s *Service) PublishToSubKey(subKey *common.Address, message Message) error {
	return s.Publish(strings.ToLower(fmt.Sprintf(SUBKEY_KEY_FMT, subKey.String())), message)
}

func (s *Service) Close(message string, p ...*MatchResponse) {
	for _, r := range p {
		s.PublishToSubKey(r.SubKey, NewTerminateMessage(message))
	}
	s.RemoveFromWaitingPool(p...)
}

// Session has been verified so begin match and send message to each player
func (srv *Service) BeginVerifiedMatch(sess *Session) error {
	if !sess.IsVerified() {
		return nil
	}
	msg, err := srv.BuildMatchVerifiedMessageWithSignature(sess)
	if err != nil {
		return err
	}
	msg.PlayerIndex = sess.Player1.Index
	msg.SignatureOpponentSubkey = sess.Player2.Token.SubKeySignature
	sess.Signature = msg.SignatureMatchHash
	payloadJson, err := util.Jsonify(msg)
	if err != nil {
		return err
	}
	relaymsg := &Message{
		Meta: &Meta{
			Code:   MATCH_VERIFIED,
			Index:  sess.Player1.Index,
			SubKey: sess.Player1.SubKey,
		},
		Payload: payloadJson,
	}
	err = srv.PublishToSubKey(relaymsg.SubKey, *relaymsg)
	if err != nil {
		return err
	}
	msg.PlayerIndex = sess.Player2.Index
	msg.SignatureOpponentSubkey = sess.Player1.Token.SubKeySignature
	payloadJson, err = util.Jsonify(msg)
	if err != nil {
		return err
	}
	relaymsg = &Message{
		Meta: &Meta{
			Code:   MATCH_VERIFIED,
			Index:  sess.Player2.Index,
			SubKey: sess.Player2.SubKey,
		},
		Payload: payloadJson,
	}
	err = srv.PublishToSubKey(relaymsg.SubKey, *relaymsg)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) NewKeyedTransactor() *bind.TransactOpts {
	privkey := s.PrivKey()
	return bind.NewKeyedTransactor(privkey)
}

func (srv *Service) BuildMatchVerifiedMessageWithSignature(s *Session) (*arcadeum.MatchVerifiedMessage, error) {
	msg := &arcadeum.MatchVerifiedMessage{
		Accounts:    [2]common.Address{*s.Player1.Account, *s.Player2.Account},
		Subkeys:     [2]common.Address{*s.Player1.SubKey, *s.Player2.SubKey},
		GameAddress: srv.ArcClient.GameAddress[s.GameID],
		Timestamp:   s.Timestamp,
		Players: [2]*arcadeum.MatchVerifiedPlayerInfo{
			{
				SeedRating:         s.Player1.Rank,
				PublicSeed:         s.Player1.SeedHash,
				SignatureTimestamp: s.Player1.TimestampSig,
			},
			{
				SeedRating:         s.Player2.Rank,
				PublicSeed:         s.Player2.SeedHash,
				SignatureTimestamp: s.Player2.TimestampSig,
			},
		},
	}
	hash, err := srv.ArcClient.MatchHash(msg)
	if err != nil {
		return nil, err
	}
	msg.MatchHash = hash

	// Have the matcher sign
	sig, err := gethcrypto.Sign(hash[:], srv.PrivKey())
	if err != nil {
		return nil, err
	}
	msg.SignatureMatchHash = &crypto.Signature{V: 27 + sig[64]}
	copy(msg.SignatureMatchHash.R[:], sig[0:32])
	copy(msg.SignatureMatchHash.S[:], sig[32:64])

	return msg, nil
}

func (s *Service) SendTimestampProof(p *PlayerInfo, timestamp int64) error {
	if p == nil {
		return nil
	}
	message := Message{
		Meta: &Meta{
			Code:   INIT,
			SubKey: p.SubKey,
			Index:  p.Index,
		},
		Payload: strconv.FormatInt(timestamp, 10)}
	err := s.PublishToSubKey(message.SubKey, message)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) RequestTimestampProof(sess *Session) error {
	log.Println("Requesting timestamp proof from both players")
	err := s.SendTimestampProof(sess.Player1, sess.Timestamp)
	if err != nil {
		return err
	}
	err = s.SendTimestampProof(sess.Player2, sess.Timestamp)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) CreateSession(p *MatchResponse) (*Session, error) {
	id := uuid.New()
	player, err := s.BuildPlayerInfo(p)
	if err != nil {
		return nil, err
	}
	return &Session{
		ID:      UUID(id),
		GameID:  player.GameID,
		Player1: player,
	}, nil
}

func (srv *Service) SignElliptic(inputs ...interface{}) (r, s *big.Int, err error) {
	compact, err := Compact(inputs...)
	if err != nil {
		return
	}
	hash := gethcrypto.Keccak256(compact)
	privkey := srv.PrivKey()
	r, s, err = ecdsa.Sign(rand.Reader, privkey, hash)
	return
}

func (srv *Service) PrivKey() *ecdsa.PrivateKey {
	path := fmt.Sprintf("%s/%s", srv.ENV.WorkingDir, srv.Config.PrivKeyFile)
	privkey, err := gethcrypto.LoadECDSA(path)
	if err != nil {
		log.Fatalf("Invalid private key: %v", err)
	}
	return privkey
}

func (srv *Service) Sign(inputs ...interface{}) ([]byte, error) {
	r, s, err := srv.SignElliptic(inputs...)
	if err != nil {
		return nil, err
	}
	return asn1.Marshal(crypto.EcdsaSignature{r, s})
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
	seedHash, err := s.ArcClient.PublicSeed(p.GameID, p.Seed)
	if err != nil {
		return nil, err
	}
	return &PlayerInfo{
		Rank:     p.Rank,
		Token:    p.Token,
		SeedHash: seedHash,
		Account:  &p.Account,
	}, nil
}

func (s *Service) OnMessage(msg *Message) error {
	sess, err := s.GetSessionBySubKey(msg.SubKey)
	if err != nil {
		return err
	}
	if sess.IsEmpty() {
		return errors.New(fmt.Sprintf("Unknown session for subkey %s", msg.SubKey.String()))
	}
	if msg.Code == SIGNED_TIMESTAMP { // Verified signed timestamp
		player := sess.FindPlayerBySubKey(msg.SubKey)
		if player == nil {
			return errors.New("Could not find player, unknown subkey")
		}
		req := &arcadeum.VerifyTimestampRequest{}
		err := json.Unmarshal([]byte(msg.Payload), req)
		if err != nil {
			return err
		}
		verified, err := s.VerifyTimestamp(req, player)
		if err != nil {
			return err
		}
		if !verified {
			return errors.New("Invalid timestamp signature proof.")
		}
		player.Verified = verified
		player.TimestampSig = req.Signature // set the verified signature
		err = s.BeginVerifiedMatch(sess)
		if err != nil {
			return err
		}
		err = s.UpdateSession(sess)
		if err != nil {
			return err
		}
	} else if !sess.IsVerified() {
		return errors.New("Match session not verified")
	} else { // verified, relay message to opponent
		opponent := sess.FindOpponent(msg.SubKey)
		if opponent == nil {
			log.Println("No opponent, swallowing message")
			return nil
		}
		s.PublishToSubKey(opponent.SubKey, *msg)
	}
	return nil
}

// Subscribed to published messages for this player
func (s *Service) SubscribeToSubKey(subkey *common.Address, messages chan *Message) {
	s.Subscribe(strings.ToLower(fmt.Sprintf(SUBKEY_KEY_FMT, subkey.String())), messages)
}
