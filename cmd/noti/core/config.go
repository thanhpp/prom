package core

import (
	"errors"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type WebServerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

var mainConfig = new(MainConfig)

func readConfigFromFile(filePath string) (err error) {
	v := viper.New()
	filePart := strings.Split(filePath, ".")
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

	return nil
}

func readConfigFromENV() {
	dbHost := os.Getenv("DBHOST")
	if len(dbHost) > 0 {
		mainConfig.DB.Host = dbHost
	}

	webHost := os.Getenv("WEBHOST")
	if len(webHost) > 0 {
		mainConfig.Web.Host = webHost
	}

	webPort := os.Getenv("WEBPORT")
	if len(webPort) > 0 {
		mainConfig.Web.Port = webPort
	}

	rabbitMQURL := os.Getenv("RABBITMQURL")
	if len(rabbitMQURL) > 0 {
		mainConfig.RabbitMQURL = rabbitMQURL
	}
}
