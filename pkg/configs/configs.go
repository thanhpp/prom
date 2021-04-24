// configs store all global config struct
package configs

type GRPCConfig struct {
	PublicHost string `mapstructure:"publichost"`
	Port       string `mapstructure:"port"`
	MaxStream  int    `mapstructure:"maxconnections"`
}

type DBConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	User string `mapstructure:"user"`
	Pass string `mapstructure:"pass"`
	Log  string `mapstructure:"log"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
	Color bool   `mapstructure:"color"`
}
