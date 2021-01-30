package controller

import (
	"HypertubeAuth/errors"
)

type Config struct {
	OauthRedirect string `conf:"oauthRedirect"`
	ServerPasswd  string `conf:"passwd"`
	ServerPort    uint   `conf:"port"`
}

var cfg *Config

func GetConfig() *Config {
	if cfg == nil {
		cfg = &Config{}
	}
	return cfg
}

func getConfig() (*Config, *errors.Error) {
	if cfg == nil {
		return nil, errors.NotConfiguredPackage.SetArgs("controller", "controller")
	}
	return cfg, nil
}

func GetServerPort() (uint, *errors.Error) {
	conf, Err := getConfig()
	if Err != nil {
		return 0, Err
	}
	return conf.ServerPort, nil
}
