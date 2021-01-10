package configsender

import "github.com/BurntSushi/toml"

type Config struct {
	Logger LoggerConf
	AMQP   AMQPConf
}

type LoggerConf struct {
	Level string
	Path  string
}

type AMQPConf struct {
	URI      string
	Qname    string
	Exchname string
	Exchtype string
}

func NewConfig(filePath string) (*Config, error) {
	var conf Config
	if _, err := toml.DecodeFile(filePath, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}
