package arcadeum

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	"github.com/horizon-games/arcadeum/server/config"
	"github.com/horizon-games/arcadeum/server/services/crypto"
	"math/big"
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
	GameID      uint32            `json:"gameID"`
	MatchID     uint32            `json:"matchID"`
	Accounts    [2]common.Address `json:"accounts"`
	Timestamp   int64             `json:"timestamp"`
	Seeds       [2][]byte         `json:"seeds"`
	SeedHashes  [2][]byte         `json:"seedhashes"`
	SeedRatings [2]uint32         `json:"seedRatings"`
	MatchHash   [32]byte          `json:"matchHash"` // signature of all fields above

	SignatureMatchHash     []byte            `json:"signature"`
	SignatureMatchHashPair *crypto.Signature `json:"signature_pair"`
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
	contract, _ := NewArcadeum(c.ArcadeumContractAddress, c.Conn)
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
	contract, err := NewArcadeum(c.ArcadeumContractAddress, c.Conn)
	if err != nil {
		return common.Address{}, err
	}

	var r1, s1, r2, s2 [32]byte
	copy(r1[:], req.Signature.R)
	copy(s1[:], req.Signature.S)
	copy(r2[:], subkeySig.R)
	copy(s2[:], subkeySig.S)
	return contract.ArcadeumCaller.PlayerAccount(
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
	contract, err := NewArcadeum(c.ArcadeumContractAddress, c.Conn)
	if err != nil {
		return common.Address{}, err
	}
	var r, s [32]byte
	copy(r[:], sig.R)
	copy(s[:], sig.S)
	return contract.ArcadeumCaller.SubkeyParent(
		&bind.CallOpts{},
		subkey,
		sig.V,
		r,
		s,
	)
}

// Check if the account owns the secret seed "deck"
func (c *Client) IsSecretSeedValid(gameID uint32, account common.Address, secretSeed []byte) (bool, error) {
	contract, err := c.DGameContract(gameID)
	if err != nil {
		return false, err
	}
	return contract.DGameCaller.IsSecretSeedValid(
		&bind.CallOpts{},
		account,
		secretSeed,
	)
}

// Calculate the balance of a given address on the Arcadeum staking contract.
func (c *Client) StakeBalanceInEth(address common.Address) (*big.Int, error) {
	contract, err := NewArcadeum(c.ArcadeumContractAddress, c.Conn)
	if err != nil {
		return nil, err
	}
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
	copy(seed1[:], msg.SeedHashes[0])
	copy(seed2[:], msg.SeedHashes[1])
	publicSeeds = [2][1][32]byte{{seed1}, {seed2}}
	gameaddr := c.GameAddress[msg.GameID]
	return contract.MatchHash(
		&bind.CallOpts{},
		gameaddr,
		msg.MatchID,
		big.NewInt(msg.Timestamp),
		msg.Accounts,
		msg.SeedRatings,
		publicSeeds,
	)

}
