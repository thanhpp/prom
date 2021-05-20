// configs store all global config struct
package configs

import "fmt"

type GRPCConfig struct {
	Name       string `mapstructure:"name"`
	PublicHost string `mapstructure:"publichost"`
	Port       string `mapstructure:"port"`
	MaxStream  int    `mapstructure:"maxconnections"`
}

type DBConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	User string `mapstructure:"user"`
	Pass string `mapstructure:"pass"`
	Name string `mapstructure:"name"`
	Log  string `mapstructure:"log"`
}

func (DB *DBConfig) GenDBDSN() (dsn string) {
	dsn = fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		DB.User, DB.Pass, DB.Name, DB.Host, DB.Port)
	return
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
	Color bool   `mapstructure:"color"`
}
