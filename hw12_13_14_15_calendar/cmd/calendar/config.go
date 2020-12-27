package main

import "github.com/BurntSushi/toml"

type Config struct {
	Logger  LoggerConf
	Storage StorageConf
	Server  ServerConf
}

type LoggerConf struct {
	Level string
	Path  string
}

type StorageConf struct {
	Type string
	Dsn  string
}

type ServerConf struct {
	Host string
	Port string
}

func NewConfig(filePath string) (*Config, error) {
	var conf Config
	if _, err := toml.DecodeFile(filePath, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}
