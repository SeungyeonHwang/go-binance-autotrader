package config

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type Config struct {
	MasterAPIKey    string
	MasterSecretKey string
	Sub1APIKey      string
	Sub1SecretKey   string
	Sub2APIKey      string
	Sub2SecretKey   string
	Sub3APIKey      string
	Sub3SecretKey   string
}

type SSMConfigLoader struct {
	SSMClient *ssm.SSM
}

func (loader *SSMConfigLoader) GetParameter(paramName string) (string, error) {
	input := &ssm.GetParameterInput{
		Name:           aws.String(paramName),
		WithDecryption: aws.Bool(true),
	}
	result, err := loader.SSMClient.GetParameter(input)
	if err != nil {
		return "", err
	}

	return *result.Parameter.Value, nil
}

func LoadConfigurationFromSSM(ssmClient *ssm.SSM) (*Config, error) {
	loader := SSMConfigLoader{SSMClient: ssmClient}
	return loadConfigFromSSM(loader)
}

func loadConfigFromSSM(loader SSMConfigLoader) (*Config, error) {
	config := &Config{}
	accounts := []string{"master", "sub1", "sub2", "sub3"}

	for _, account := range accounts {
		apiKeyName := fmt.Sprintf("/binance/%s/api_key", account)
		secretKeyName := fmt.Sprintf("/binance/%s/secret_key", account)

		apiKey, err := loader.GetParameter(apiKeyName)
		if err != nil {
			return nil, err
		}

		secretKey, err := loader.GetParameter(secretKeyName)
		if err != nil {
			return nil, err
		}

		switch account {
		case "master":
			config.MasterAPIKey = apiKey
			config.MasterSecretKey = secretKey
		case "sub1":
			config.Sub1APIKey = apiKey
			config.Sub1SecretKey = secretKey
		case "sub2":
			config.Sub2APIKey = apiKey
			config.Sub2SecretKey = secretKey
		case "sub3":
			config.Sub3APIKey = apiKey
			config.Sub3SecretKey = secretKey
		}
	}

	return config, nil
}
