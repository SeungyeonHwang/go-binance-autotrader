package binance

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Binance map[string]struct {
		APIKey    string `yaml:"api_key"`
		SecretKey string `yaml:"secret_key"`
	} `yaml:"binance"`
}

func getConfig() (Config, error) {
	data, err := ioutil.ReadFile("config/settings.yaml")
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
