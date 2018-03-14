package arcadeum

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	"github.com/horizon-games/arcadeum/server/config"
	"github.com/horizon-games/arcadeum/server/services/crypto"
	"github.com/horizon-games/arcadeum/server/services/util"
)

const (
	PRICE_SAMPLE_INTERVAL_IN_SEC = 60 // 1 min
)

type EtherScanPriceResponse struct {
	Status  int         `json:"status,string"`
	Message string      `json:"message"`
	Result  PriceResult `json:"result"`
}

type PriceResult struct {
	EthBTC          float64 `json:"ethbtc,string"`           // price of 1 ETH in BTC
	EthBTCTUnixTime int64   `json:"ethbtc_timestamp,string"` // unix time of last Eth/BTC price calculation
	EthUSD          float64 `json:"ethusd,string"`           // price of 1 ETH in USD
	EthUSDUnixTime  int64   `json:"ethusd_timestamp,string"` // unix time of last Eth/USD price calculation
}

type VerifyTimestampRequest struct {
	MatchID   uint32           `json:"matchID"`
	GameID    uint32           `json:"gameID"`
	Timestamp int64            `json:"timestamp"` // unix time
	Signature crypto.Signature `json:"signature"` // as signed by the players private key
}

type Client struct {
	LastPriceUnixTime int64   // time of last ether/usd price retrieval; used for cache refresh
	USDInETH          float64 // price of 1 USD in ETH
	EtherscanPriceURL string

	MinStakeUSD             float32 // minimum Arcadeum stake amount in USD
	ArcadeumContractAddress common.Address
	GameAddress             map[uint32]common.Address // gameID -> game contract address
	Conn                    *ethclient.Client

	ArcadeumContract *Arcadeum
}

type MatchVerifiedMessage struct {
	Accounts [2]common.Address

	GameAddress common.Address              `json:"game"`
	MatchID     uint32                      `json:"matchID"`
	Timestamp   int64                       `json:"timestamp"`
	PlayerIndex uint8                       `json:"playerID"`
	Players     [2]*MatchVerifiedPlayerInfo `json:"players"`

	MatchHash                  [32]byte
	SignatureMatchHash         *crypto.Signature `json:"matchSignature"`
	SignatureOpponentTimestamp *crypto.Signature `json:"opponentTimestampSignature"`
	SignatureOpponentSubkey    *crypto.Signature `json:"opponentSubkeySignature"`
}

type MatchVerifiedPlayerInfo struct {
	SeedRating uint32 `json:"seedRating"`
	PublicSeed []byte `json:"publicSeed,string"`
}

type StakedStatus int

const (
	UNKNOWN                     = -1
	NOT_STAKED                  = 0
	STAKED                      = 1
	STAKED_INSUFFICIENT_BALANCE = 2 // player's stake balance is below latest calculated stake
)

type IWithdrawalStartedHandler interface {
	OnWithdrawalStarted(event *ArcadeumWithdrawalStarted)
}

var withdrawalStartedChan = make(chan *ArcadeumWithdrawalStarted)
var withdrawalStartedSubscription event.Subscription

func (c *Client) SubscribeWithdrawalStarted() (event.Subscription, error) {
	contract := c.ArcadeumContract
	sub, err := contract.ArcadeumFilterer.WatchWithdrawalStarted(
		&bind.WatchOpts{},
		withdrawalStartedChan,
		[]common.Address{c.ArcadeumContractAddress},
	)
	if err != nil {
		return nil, err
	}
	withdrawalStartedSubscription = sub
	return sub, err
}

func (c *Client) HandleWithdrawalStarted(handler IWithdrawalStartedHandler) {
	for {
		ev := <-withdrawalStartedChan
		handler.OnWithdrawalStarted(ev)
	}
}

func (c *Client) VerifySignedTimestamp(
	gameID uint32,
	matchID uint32,
	req *VerifyTimestampRequest,
	subkeySig crypto.Signature) (common.Address, error) {
	gameaddr := c.GameAddress[gameID]
	contract := c.ArcadeumContract
	sigR, err := util.DecodeHexString(string(req.Signature.R))
	if err != nil {
		return common.Address{}, err
	}
	sigS, err := util.DecodeHexString(string(req.Signature.S))
	if err != nil {
		return common.Address{}, err
	}
	subkeyR, err := util.DecodeHexString(string(subkeySig.R))
	if err != nil {
		return common.Address{}, err
	}
	subkeyS, err := util.DecodeHexString(string(subkeySig.S))
	if err != nil {
		return common.Address{}, err
	}
	var r1, s1, r2, s2 [32]byte
	copy(r1[:], sigR)
	copy(s1[:], sigS)
	copy(r2[:], subkeyR)
	copy(s2[:], subkeyS)
	return contract.ArcadeumCaller.PlayerAccountXXX(
		&bind.CallOpts{},
		gameaddr,
		matchID,
		big.NewInt(req.Timestamp),
		req.Signature.V,
		r1,
		s1,
		subkeySig.V,
		r2,
		s2,
	)
}

// Return the price of 1 USD in ETH
func (c *Client) PriceUSDInEth() (float64, error) {
	currentTime := time.Now().Unix()
	if currentTime-c.LastPriceUnixTime > PRICE_SAMPLE_INTERVAL_IN_SEC {
		resp, err := http.Get(c.EtherscanPriceURL)
		if err != nil {
			if c.USDInETH > 0 {
				log.Println("error requesting price from etherscan returning cached price", c.USDInETH)
				return c.USDInETH, nil
			}
			return 0, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			if c.USDInETH > 0 {
				log.Println("Non-200 response from etherscan returning cached price", c.USDInETH)
				return c.USDInETH, nil
			}
			return 0, errors.New(fmt.Sprintf("Etherscan status code %d", resp.StatusCode))
		}
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			if c.USDInETH > 0 {
				log.Println("Error reading etherscan response body, returning cached price", c.USDInETH)
				return c.USDInETH, nil // return the last price
			}
			return 0, err
		}
		r := &EtherScanPriceResponse{}
		err2 := json.Unmarshal(bodyBytes, r)
		if err2 != nil {
			if c.USDInETH > 0 {
				log.Println("Error deserializing etherscan response body, returning cached price", c.USDInETH)
				return c.USDInETH, nil // return the last price
			}
			return 0, err
		}
		c.LastPriceUnixTime = currentTime
		c.USDInETH = 1 / r.Result.EthUSD
	}
	return c.USDInETH, nil
}

func NewArcadeumClient(ethcfg *config.ETHConfig, arccfg *config.ArcadeumConfig) *Client {
	if !common.IsHexAddress(arccfg.ContractAddress) {
		log.Fatalf("Invalid Arcadeum contract address %s.", arccfg.ContractAddress)
	}
	arc := common.HexToAddress(arccfg.ContractAddress)
	conn, err := ethclient.Dial(ethcfg.NodeURL)
	if err != nil {
		log.Fatalf("Could not create client RPC to node ", ethcfg.NodeURL)
	}
	address := make(map[uint32]common.Address)
	for _, game := range arccfg.Games {
		if !common.IsHexAddress(game.ContractAddress) {
			log.Fatalf("Invalid game contract address %s.", game.ContractAddress)
		}
		address[game.ID] = common.HexToAddress(game.ContractAddress)
	}
	contract, err := NewArcadeum(arc, conn)
	if err != nil {
		log.Fatalf("Failure loading Arcadeum contract.", err)
	}
	client := &Client{
		MinStakeUSD:             arccfg.MinStakeUSD,
		ArcadeumContractAddress: arc,
		Conn:              conn,
		GameAddress:       address,
		EtherscanPriceURL: arccfg.EtherscanPriceURL,
		ArcadeumContract:  contract,
	}
	client.SubscribeWithdrawalStarted()
	return client
}

func (c *Client) MatchDuration(gameID uint32) (time.Duration, error) {
	contract, err := c.DGameContract(gameID)
	if err != nil {
		return 0, err
	}
	duration, err := contract.MatchDuration(&bind.CallOpts{})
	if err != nil {
		return 0, err
	}
	return time.Duration(duration.Int64()) * time.Second, nil
}

// Call constant function ETH contract, passing in the payload, address of requester.
// Use gameID to map to correct game contract.
func (c *Client) CalculateRank(gameID uint32, secretSeed []byte) (uint32, error) {
	contract, err := c.DGameContract(gameID)
	if err != nil {
		return 0, err
	}
	return contract.DGameCaller.SecretSeedRating(
		&bind.CallOpts{},
		secretSeed,
	)
}

// Return the address of the signer of the subkey
func (c *Client) SubKeyParent(subkey common.Address, sig crypto.Signature) (common.Address, error) {
	contract := c.ArcadeumContract

	r, err := util.DecodeHexString(string(sig.R))
	if err != nil {
		return common.BytesToAddress(make([]byte, 20)), err
	}
	s, err := util.DecodeHexString(string(sig.S))
	if err != nil {
		return common.BytesToAddress(make([]byte, 20)), err
	}

	var r32 [32]byte
	var s32 [32]byte
	copy(r32[:], r)
	copy(s32[:], s)

	return contract.ArcadeumCaller.SubkeyParentXXX(
		&bind.CallOpts{},
		subkey,
		sig.V,
		r32,
		s32,
	)
}

// Return the subkey that signed a given timestamp
func (c *Client) TimestampSubkey(timestamp int64, sig crypto.Signature) (common.Address, error) {
	contract := c.ArcadeumContract

	r, err := util.DecodeHexString(string(sig.R))
	if err != nil {
		return common.BytesToAddress(make([]byte, 20)), err
	}
	s, err := util.DecodeHexString(string(sig.S))
	if err != nil {
		return common.BytesToAddress(make([]byte, 20)), err
	}

	var r32 [32]byte
	var s32 [32]byte
	copy(r32[:], r)
	copy(s32[:], s)

	return contract.ArcadeumCaller.TimestampSubkeyXXX(
		&bind.CallOpts{},
		big.NewInt(timestamp),
		sig.V,
		r32,
		s32,
	)
}

// Check if the account owns the secret seed "deck"
func (c *Client) IsSecretSeedValid(gameID uint32, account common.Address, secretSeed []byte) (bool, error) {
	contract, err := c.DGameContract(gameID)
	if err != nil {
		return false, err
	}
	ss, err := util.DecodeHexString(string(secretSeed))
	if err != nil {
		return false, err
	}
	return contract.DGameCaller.IsSecretSeedValid(
		&bind.CallOpts{},
		account,
		ss,
	)
}

// Calculate the balance of a given address on the Arcadeum staking contract.
func (c *Client) StakeBalanceInEth(address common.Address) (*big.Int, error) {
	contract := c.ArcadeumContract
	return contract.ArcadeumCaller.Balance(&bind.CallOpts{}, address)
}

func (c *Client) GetStakedStatus(from common.Address) (StakedStatus, error) {
	price, err := c.PriceUSDInEth()
	if err != nil {
		return UNKNOWN, err
	}
	balance, err := c.StakeBalanceInEth(from)
	if err != nil {
		return UNKNOWN, err
	}
	log.Println("price balance min_stake", price, balance, price*float64(c.MinStakeUSD))
	if float64(balance.Uint64()) > price*float64(c.MinStakeUSD) {
		return STAKED, nil
	} else if balance.Uint64() > 0 {
		return STAKED_INSUFFICIENT_BALANCE, nil
	}

	return NOT_STAKED, nil
}

func (c *Client) DGameContract(gameID uint32) (*DGame, error) {
	gameaddr := c.GameAddress[gameID]
	return NewDGame(gameaddr, c.Conn)
}

func (c *Client) PublicSeed(gameID uint32, secretSeed []byte) ([]byte, error) {
	contract, err := c.DGameContract(gameID)
	if err != nil {
		return nil, err
	}
	res, err := contract.DGameCaller.PublicSeed(
		&bind.CallOpts{},
		secretSeed,
	)
	if err != nil {
		return nil, err
	}
	return res[0][:], nil
}

func (c *Client) MatchHash(msg *MatchVerifiedMessage) ([32]byte, error) {
	contract := c.ArcadeumContract
	// Due to a bug in abigen, you have to hack the solidity code so abigen
	// produces something remotely workable. Hence the awful datatype [2][1][32]byte
	var publicSeeds [2][1][32]byte
	var seed1, seed2 [32]byte
	copy(seed1[:], msg.Players[0].PublicSeed)
	copy(seed2[:], msg.Players[1].PublicSeed)
	publicSeeds = [2][1][32]byte{{seed1}, {seed2}}
	return contract.MatchHash(
		&bind.CallOpts{},
		msg.GameAddress,
		msg.MatchID,
		big.NewInt(msg.Timestamp),
		msg.Accounts,
		[2]uint32{msg.Players[0].SeedRating, msg.Players[1].SeedRating},
		publicSeeds,
	)

}
