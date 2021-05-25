package core

import (
	"github.com/thanhpp/prom/pkg/configs"
)

type MainConfig struct {
	ServiceName string               `mapstructure:"servicename"`
	Environment string               `mapstructure:"environment"`
	Log         configs.LoggerConfig `mapstructure:"log"`
	DB          configs.DBConfig     `mapstructure:"database"`
	Web         WebServerConfig      `mapstructure:"webserver"`
	RabbitMQURL string               `mapstructure:"rabbitmqurl"`
}

func SetConfig(filePath string) (err error) {
	if err = readConfigFromFile(filePath); err != nil {
		return err
	}

	readConfigFromENV()
	return nil
}

func GetConfig() *MainConfig {
	return mainConfig
}
