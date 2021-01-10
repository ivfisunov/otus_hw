package configscheduler

import "github.com/BurntSushi/toml"

type Config struct {
	Logger LoggerConf
	Grpc   APIConf
	Amqp   AmqpConf
}

type LoggerConf struct {
	Level string
	Path  string
}

type APIConf struct {
	Host string
	Port string
}

type AmqpConf struct {
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
