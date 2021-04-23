package core

// --------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- CONFIG ----------------------------------------------------------

type MainConfig struct {
	DB  dbConfig     `mapstructure:"database"`
	Log loggerConfig `mapstructure:"log"`
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
