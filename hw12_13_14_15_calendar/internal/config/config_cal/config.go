package configcal

import "github.com/BurntSushi/toml"

type Config struct {
	Logger  LoggerConf
	Storage StorageConf
	HTTP    HTTPConf
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

type HTTPConf struct {
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
