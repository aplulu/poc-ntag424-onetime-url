package config

import (
	"encoding/hex"
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Listen             string   `envconfig:"listen" default:""`
	Port               string   `envconfig:"port" default:"8080"`
	Key                string   `envconfig:"key" default:"00000000000000000000000000000000"`
	CORSAllowedOrigins []string `envconfig:"cors_allowed_origins" default:"*"`
	CORSMaxAge         int      `envconfig:"cors_max_age" default:"600"`
}

var conf config
var key []byte

// LoadConf loads the configuration from the environment variables.
func LoadConf() error {
	if err := envconfig.Process("", &conf); err != nil {
		return fmt.Errorf("config.LoadConf: failed to load conf: %w", err)
	}

	var err error
	key, err = hex.DecodeString(conf.Key)
	if err != nil {
		return fmt.Errorf("config.LoadConf: failed to decode key: %w", err)
	}

	return nil
}

func Listen() string {
	return conf.Listen
}

func Port() string {
	return conf.Port
}

func Key() []byte {
	return key
}

func CORSAllowedOrigins() []string {
	return conf.CORSAllowedOrigins
}

func CORSMaxAge() int {
	return conf.CORSMaxAge
}
