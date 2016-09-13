package config

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

const configFile = "config.ini"

var cfg *ConfigWrapper

func LoadConfig() *ConfigWrapper {
	//If we don't already have a pointer to a config struct, load it from a file
	if cfg == nil {
		_, err := os.Stat(configFile)
		if err != nil {
			log.Fatal("Config file is missing: ", configFile)
		}
		if _, err := toml.DecodeFile(configFile, &cfg); err != nil {
			log.Fatal(err)
		}
	}
	return cfg
}

type appConfig struct {
	AuthToken     string
	ApplicationID string
	CommandPrefix string
}

type ConfigWrapper struct {
	AppConfig appConfig
}
