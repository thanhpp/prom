package core

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
)

type dbConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	User string `mapstructure:"user"`
	Pass string `mapstructure:"pass"`
	Log  string `mapstructure:"log"`
}

type loggerConfig struct {
	Level string `mapstructure:"level"`
	Color bool   `mapstructure:"color"`
}

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

	return nil
}