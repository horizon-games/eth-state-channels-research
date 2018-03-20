package config

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type ENVConfig struct {
	Mode        string `toml:"mode"`
	Environment string `toml:"environment"`
	DebugMode   bool   `toml:"debug_mode"`
	Port        int    `toml:"port"`
	TLSEnabled  bool   `toml:"tls_enabled"`
	TLSCertFile string `toml:"tls_cert_file"`
	TLSKeyFile  string `toml:"tls_key_file"`
	WorkingDir  string `toml:"working_dir"`
}

type MatcherConfig struct {
	AccountAddress common.Address // Ethereum address of account associated with this private key (derived)
	PrivKeyFile    string         `toml:"priv_key_file"`
}

type ETHConfig struct {
	NodeURL string `toml:"node_url"`
}

type ArcadeumConfig struct {
	EtherscanPriceURL string     `toml:"etherscan_price_url"` // endpoint to get latest ETH/USD price
	MinStakeUSD       float32    `toml:"min_stake_usd"`
	ContractAddress   string     `toml:"contract_address"` // arcadeum contract address
	Games             []GameInfo `toml:"games"`
}

type GameInfo struct {
	ID              uint32 `toml:"id"`
	Name            string `toml:"name"`
	ContractAddress string `toml:"contract_address"`
}

type RedisConfig struct {
	Address  string `toml:"address"`
	Password string `toml:"password"`
}

type Config struct {
	ENV            ENVConfig      `toml:"env"`
	MatcherConfig  MatcherConfig  `toml:"matcher"`
	ETHConfig      ETHConfig      `toml:"eth"`
	ArcadeumConfig ArcadeumConfig `toml:"arcadeum"`
	RedisConfig    RedisConfig    `toml:"redis"`
}

func NewFromFile(file string, env string, config interface{}) error {
	if file == "" {
		file = env
	}

	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return errors.Wrap(err, "failed to load config file")
	}

	if _, err := toml.DecodeFile(file, config); err != nil {
		return errors.Wrap(err, "failed to parse config file")
	}

	return nil
}
