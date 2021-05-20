package core

import (
	"github.com/thanhpp/prom/pkg/configs"
	"github.com/thanhpp/prom/pkg/etcdclient"
)

// --------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- CONFIG ----------------------------------------------------------

type MainConfig struct {
	DB      configs.DBConfig       `mapstructure:"database"`
	Log     configs.LoggerConfig   `mapstructure:"log"`
	ETCD    etcdclient.ETCDConfigs `mapstructure:"etcd"`
	GRPC    configs.GRPCConfig     `mapstructure:"grpc"`
	ShardID int64                  `mapstructure:"shardid"`
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
