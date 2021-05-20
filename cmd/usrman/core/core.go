package core

import (
	"github.com/thanhpp/prom/pkg/configs"
	"github.com/thanhpp/prom/pkg/etcdclient"
)

type MainConfig struct {
	DB   configs.DBConfig       `mapstructure:"database"`
	Log  configs.LoggerConfig   `mapstructure:"log"`
	GRPC configs.GRPCConfig     `mapstructure:"grpc"`
	ETCD etcdclient.ETCDConfigs `mapstructure:"etcd"`
}

var mainConfig = new(MainConfig)

func GetConfig() *MainConfig {
	return mainConfig
}

func SetMainConfig(configPath string) (err error) {
	if err := readConfigFromFile(configPath); err != nil {
		return err
	}

	return nil
}
