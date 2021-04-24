package core

import "github.com/thanhpp/prom/pkg/configs"

// --------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- CONFIG ----------------------------------------------------------

type MainConfig struct {
	DB   configs.DBConfig     `mapstructure:"database"`
	Log  configs.LoggerConfig `mapstructure:"log"`
	GRPC configs.GRPCConfig   `mapstructure:"grpc"`
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
