package matcher

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-redis/redis"
	"github.com/horizon-games/arcadeum/server/config"
	"github.com/horizon-games/arcadeum/server/services/crypto"

	"encoding/json"
	"github.com/horizon-games/arcadeum/server/services/util"
	"github.com/pkg/errors"
	"log"
	"time"
)

type UUID string

const (
	RANK_KEY_FMT = "rank:%d"
)

type PlayerInfo struct {
	Rank  uint32 `json:"rank"`  // Player rank as returned from Solidity pure function
	Index uint8  `json:"index"` // index in match session; arbitrarily set when match found and player joins a session
	*Token

	SeedHash     []byte            `json:"seedHash,string"` // Hash of seed as returned from Solidity pure function
	Account      *common.Address   `json:"account,string"`  // owner account of signed subkey (derived from Token); See Token
	TimestampSig *crypto.Signature `json:"timestampSig"`    // TimestampSig of session timestamp by this player
	Verified     bool              `json:"verified"`        // true if the player has proven their TimestampSig
}

// Represents a matched game session
type Session struct {
	ID        UUID              `json:"id,string"`
	GameID    uint32            `json:"gameID"`
	Player1   *PlayerInfo       `json:"player1"`
	Player2   *PlayerInfo       `json:"player2"`
	Timestamp int64             `json:"timestamp"` // Game start in Unix time
	Signature *crypto.Signature `json:"signature"` // Matcher's session signature
}

type SessionManager struct {
	SessionPool map[UUID]*Session // map session UUID -> in-game session
	RedisClient *redis.Client
}

func (s *Session) UUID() string {
	return string(s.ID)
}

func (s *Session) IsEmpty() bool {
	return string(s.ID) == ""
}

func (s *Session) FindPlayerByAccount(account common.Address) (*PlayerInfo, error) {
	if s.Player1.Account != nil && *s.Player1.Account == account {
		return s.Player1, nil
	} else if s.Player2.Account != nil && *s.Player2.Account == account {
		return s.Player2, nil
	}
	return nil, errors.New("Unknonw account " + account.String())
}

func (u UUID) IsEmpty() bool {
	return string(u) == ""
}

func NewSessionManager(rediscfg *config.RedisConfig) *SessionManager {
	redis := redis.NewClient(&redis.Options{
		Addr:     rediscfg.Address,
		Password: rediscfg.Password,
		DB:       0, // use default DB
	})
	_, err := redis.Ping().Result()
	if err != nil {
		log.Fatalf("Redis server unreachable: %s", err.Error())
	}
	return &SessionManager{
		SessionPool: make(map[UUID]*Session),
		RedisClient: redis,
	}
}

// Find a similarly ranked player to play, pseudo-randomly chosen.
// Returns the UUID of the game session. This method will remove the waiting
// session from the wait pool. If a session is not found and there are no errors,
// an empty UUID is returned.
func (mgr *SessionManager) TakeRandomSessionByRank(rank uint32) (UUID, error) {
	uuid, err := mgr.RedisClient.SPop(fmt.Sprintf(RANK_KEY_FMT, rank)).Result()
	if err == redis.Nil { // key not found
		return UUID(""), nil
	}
	if err != nil {
		return UUID(""), err
	}
	return UUID(uuid), nil
}

func (mgr *SessionManager) ReaddToMatchPool(rank uint32, uid UUID) error {
	count, err := mgr.RedisClient.SAdd(fmt.Sprintf(RANK_KEY_FMT, rank), uid).Result() // add session to wait pool set
	if err != nil {
		return err
	}
	if count < 1 {
		return errors.New("Failed to re-add session to match pool")
	}
	return nil
}

func (mgr *SessionManager) AddToMatchPool(sess *Session) error {
	err := mgr.RedisClient.Watch(mgr.addSessionToWaitPoolTx(sess.Player1.Rank, sess))
	if err != nil {
		return err
	}
	return nil
}

func (mgr *SessionManager) addSessionToWaitPoolTx(rank uint32, sess *Session) func(t *redis.Tx) error {
	return func(tx *redis.Tx) error {
		uid := sess.UUID()
		err := setSessionKeys(tx, sess)
		if err != nil {
			return err
		}
		count, err := tx.SAdd(fmt.Sprintf(RANK_KEY_FMT, rank), uid).Result() // add session to wait pool set
		if err != nil {
			return err
		}
		if count < 1 {
			return errors.New("Unable to add session to waitpool")
		}
		return nil
	}
}

func (mgr *SessionManager) UpdateSession(sess *Session) error {
	err := mgr.RedisClient.Watch(mgr.updateSessionTx(sess))
	if err != nil {
		return err
	}
	return nil
}

func (mgr *SessionManager) updateSessionTx(sess *Session) func(t *redis.Tx) error {
	return func(tx *redis.Tx) error {
		return setSessionKeys(tx, sess)
	}
}

// There are three keys stored in redis for searchability conveninence.
// 1) uid -> session hash (key/value pairs of all session data)
// 2) player account -> uid
// 3) player subkey -> uid
// These keys are needed in order to search for session info given one of uid, account or subkey. In the case
// of account and subkey keys, two redis calls need to be made to fetch the session.
func setSessionKeys(tx *redis.Tx, sess *Session) error {
	uid := sess.UUID()
	sessJson, err := util.Jsonify(sess)
	if err != nil {
		return err
	}
	setSessStatus, err := tx.Set(uid, sessJson, time.Hour).Result() // overrides previous
	if err != nil {
		return err
	}
	if setSessStatus == "" {
		return errors.New("Error updating sesssion data")
	}
	err = setPlayerSessionKeys(tx, uid, sess.Player1)
	if err != nil {
		return err
	}
	err = setPlayerSessionKeys(tx, uid, sess.Player2)
	if err != nil {
		return err
	}
	return nil
}

func setPlayerSessionKeys(tx *redis.Tx, uid string, p *PlayerInfo) error {
	if p != nil {
		status, err := tx.Set(p.SubKey.String(), uid, time.Hour).Result()
		if err != nil {
			return err
		}
		if status == "" {
			return errors.New(fmt.Sprintf("Error setting player subkey %s", p.SubKey.String()))
		}
		status, err = tx.Set(p.Account.String(), uid, time.Hour).Result()
		if err != nil {
			return err
		}
		if status == "" {
			return errors.New(fmt.Sprintf("Error setting account key %s", p.Account.String()))
		}
	}
	return nil
}

func (mgr *SessionManager) GetSessionByID(uid UUID) (*Session, error) {
	sess := mgr.SessionPool[uid]
	if sess != nil {
		return sess, nil
	}
	sess = &Session{}
	jsonStr, err := mgr.RedisClient.Get(string(uid)).Result()
	if err == redis.Nil {
		return sess, nil
	}
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(jsonStr), sess)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func (mgr *SessionManager) GetSessionBySubKey(key *common.Address) (*Session, error) {
	return mgr.getSessionByKey(key)
}

func (mgr *SessionManager) GetSessionByAccount(key *common.Address) (*Session, error) {
	return mgr.getSessionByKey(key)
}

func (mgr *SessionManager) getSessionByKey(key *common.Address) (*Session, error) {
	uid, err := mgr.RedisClient.Get(key.String()).Result()
	if err == redis.Nil {
		return &Session{}, nil
	}
	if err != nil {
		return nil, err
	}
	return mgr.GetSessionByID(UUID(uid))

}
