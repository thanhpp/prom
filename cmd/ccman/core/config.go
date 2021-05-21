package core

import (
	"errors"
	"os"
	"strings"

	"github.com/spf13/viper"
)

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

func setConfigFromENV() {
	etcdEndpoints := os.Getenv("ETCDENDPOINT")
	if len(etcdEndpoints) > 0 {
		strs := strings.Split(etcdEndpoints, ",")
		mainConfig.ETCD.Endpoints = strs
	}

	dbHost := os.Getenv("DBHOST")
	if len(dbHost) > 0 {
		mainConfig.DB.Host = dbHost
	}

	dockerMode := os.Getenv("DOCKERMODE")
	if strings.EqualFold(dockerMode, "true") {
		mainConfig.GRPC.DockerMode = true
	}
}
