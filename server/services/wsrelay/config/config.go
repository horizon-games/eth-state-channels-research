package config

import (
	"github.com/horizon-games/dgame-server/config"
)

type Config struct {
	ENV            config.ENVConfig      `toml:"env"`
	MatcherConfig  config.MatcherConfig  `toml:"matcher"`
	ETHConfig      config.ETHConfig      `toml:"eth"`
	ArcadeumConfig config.ArcadeumConfig `toml:"arcadeum"`
}
