package config

import (
	"fmt"
	logging "github.com/ipfs/go-log/v2"
	"github.com/spf13/viper"
	"path/filepath"
)

var log = logging.Logger("config")

type Config struct {
	Network   string       `json:"network"`
	Eth1Rpc   string       `json:"eth1rpc"`
	Eth2Rpc   string       `json:"eth2rpc"`
	Store     StoreSetting `json:"store"`
	EtherScan EtherScan    `json:"etherscan"`
}

type StoreSetting struct {
	User    string `yaml:"user"`
	Pass    string `yaml:"pass"`
	Host    string `yaml:"host"`
	DB      string `yaml:"db"`
	LogMode string `yaml:"logmode"`
}

type EtherScan struct {
	Endpoint string `yaml:"endpoint"`
	ApiKey   string `yaml:"apikey"`
}

func InitConfig(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(filepath.Dir(path))

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var conf Config
	err = v.Unmarshal(&conf)
	if err != nil {
		return nil, err
	}

	if err := conf.Validate(); err != nil {
		return nil, err
	}

	return &conf, nil
}

func (cfg *Config) Validate() error {
	if cfg.Network != "holesky" && cfg.Network != "mainnet" {
		return fmt.Errorf("invalid network: %v", cfg.Network)
	}
	if cfg.Eth1Rpc == "" || cfg.Eth2Rpc == "" {
		return fmt.Errorf("invalid network: %v", cfg.Network)
	}

	if cfg.EtherScan.ApiKey == "" {
		return fmt.Errorf("invalid etherscan apiKey")
	}
	if cfg.EtherScan.Endpoint == "" {
		return fmt.Errorf("invalid etherscan endpoint")
	}

	return nil
}
