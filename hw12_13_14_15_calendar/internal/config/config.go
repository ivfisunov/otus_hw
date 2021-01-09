package config

import "github.com/BurntSushi/toml"

type Config struct {
	Logger  LoggerConf
	Storage StorageConf
	Http    HttpConf
	Grpc    GrpcConf
}

type LoggerConf struct {
	Level string
	Path  string
}

type StorageConf struct {
	Type string
	Dsn  string
}

type HttpConf struct {
	Host string
	Port string
}

type GrpcConf struct {
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
