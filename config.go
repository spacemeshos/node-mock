package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type config struct {
	RPCPort uint
}

// ConfigError config read and parse errors
type ConfigError string

func (e ConfigError) Error() string {
	return string(e)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func initConfig(fileName string) (*config, error) {
	var config config

	if !fileExists(fileName) {
		return nil, ConfigError(fmt.Sprintf("can`t find config file '%s'", fileName))
	}

	_, err := toml.DecodeFile(fileName, &config)

	return &config, err
}
