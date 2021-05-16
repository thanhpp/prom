package core

func SetMainConfig(path string) (err error) {
	if err := readConfigFromFile(path); err != nil {
		return err
	}

	return nil
}

func GetMainConfig() *MainConfig {
	return mainConfig
}
