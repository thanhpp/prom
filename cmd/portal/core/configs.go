package core

import (
	"errors"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/thanhpp/prom/cmd/portal/repository/redisdb"
	"github.com/thanhpp/prom/pkg/configs"
	"github.com/thanhpp/prom/pkg/etcdclient"
)

type MainConfig struct {
	ServiceName string                 `mapstructure:"servicename"`
	Environment string                 `mapstructure:"environment"`
	Log         configs.LoggerConfig   `mapstructure:"log"`
	ETCD        etcdclient.ETCDConfigs `mapstructure:"etcd"`
	Redis       redisdb.RedisConfig    `mapstructure:"redis"`
	WebServer   WebServerConfig        `mapstructure:"webserver"`
	RabbitMQURL string                 `mapstructure:"rabbitmqurl"`
}

type WebServerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

var mainConfig = new(MainConfig)

func readConfigFromFile(filepath string) (err error) {
	v := viper.New()
	filePart := strings.Split(filepath, ".")
	if len(filePart) != 2 {
		return errors.New("Unacceptable file path format. Require *.*")
	}
	v.SetConfigName(filePart[0])
	v.SetConfigType(filePart[1])

	// add config path
	v.AddConfigPath(".")
	v.AddConfigPath("../")
	v.AddConfigPath("../../")

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	if err := v.Unmarshal(mainConfig); err != nil {
		return err
	}

	return
}

func setConfigFromENV() {
	etcdEndpoints := os.Getenv("ETCDENDPOINT")
	if len(etcdEndpoints) > 0 {
		strs := strings.Split(etcdEndpoints, ",")
		mainConfig.ETCD.Endpoints = strs
	}

	redisAddr := os.Getenv("REDISADDR")
	if len(redisAddr) > 0 {
		mainConfig.Redis.Addr = redisAddr
	}

	webHost := os.Getenv("WEBHOST")
	if len(webHost) > 0 {
		mainConfig.WebServer.Host = webHost
	}

	webPort := os.Getenv("WEBPORT")
	if len(webPort) > 0 {
		mainConfig.WebServer.Port = webPort
	}

	rabbitMQURL := os.Getenv("RABBITMQURL")
	if len(rabbitMQURL) > 0 {
		mainConfig.RabbitMQURL = rabbitMQURL
	}
}
